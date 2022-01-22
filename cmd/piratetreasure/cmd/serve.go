/*
Copyright © 2019-2021 Netskope
*/

package cmd

import (
	"context"
	"fmt"

	"github.com/go-openapi/spec"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	health "github.com/netskope-qe/toucan-base/api/proto/toucanbase"
	toucan_handlers "github.com/netskope-qe/toucan-base/pkg/handlers"
	toucan_service "github.com/netskope-qe/toucan-base/pkg/service"
	cfg "github.com/netskope/go-kestrel/pkg/config"
	"github.com/netskope/go-kestrel/pkg/log"
	"github.com/netskope/go-kestrel/pkg/server"

	apis "github.com/netskope/piratetreasure/api/proto/piratetreasure"
	"github.com/netskope/piratetreasure/internal/build"
	"github.com/netskope/piratetreasure/internal/handlers"
)

// Constants used for serve command and related flags
const (
	ServeCmd      = "serve"
	ServeCmdShort = "Start the service"
	ServeCmdLong  = "Starts the service on the specified {host}:{port} and listens for requests"

	NoGRPCGatewayFlag      = "no-grpc-gateway"
	NoGRPCGatewayFlagShort = "g"
	NoSwaggerUIFlag        = "no-swagger-ui"
	NoSwaggerUIFlagShort   = "s"

	NoGRPCGatewayFlagUsage = "disables the REST<->gRPC gateway and the Swagger UI"
	NoSwaggerUIFlagUsage   = "disables the Swagger UI (http(s)://{host}:{port}/swagger-ui/)"

	HostFlagUsage = "the host to listen on"
	PortFlagUsage = "the port to bind to"

	HostFlag      = "host"
	PortFlag      = "port"
	PortFlagShort = "p"

	HostDefaultValue = "0.0.0.0"
	PortDefaultValue = 12345
)

// serviceHandlers struct holds the handlers for the service piratetreasure.
// TODO :: Replace HelloWorld with your own
// TODO :: Add additional handlers as needed. Make sure they are initialized (below)
// when the service starts.
type serviceHandlers struct {
	serviceHealth   toucan_handlers.HealthServiceHandler
	serviceTreasure handlers.TreasureServiceHandler
}

var (
	// the command
	serveCmd = &cobra.Command{
		Use:   ServeCmd,
		Short: ServeCmdShort,
		Long:  ServeCmdLong,
		Run:   serveCmdFunc(),
	}

	// the service logger
	logger = log.NewLogger("piratetreasure")

	// holds the service handlers
	svcHandlers *serviceHandlers
)

// panic when the logger is nil
func panicForNilLogger() {
	if logger == nil {
		panic("logger is not initialized")
	}
}

// unary middleware which invokes a callback during the request
func requestCallBackUnaryInterceptor(callback func()) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		callback()
		return handler(ctx, req)
	}
}

// stream middlware which invokes a callback during the request
func requestCallBackStreamInterceptor(callback func()) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		callback()
		return handler(srv, stream)
	}
}

// unary middleware for logging the request to stdout
func requestLoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		panicForNilLogger()

		md, _ := metadata.FromIncomingContext(ctx)
		logger.Info("request from caller",
			zap.String("info.full_method", info.FullMethod),
			zap.Any("request", req),
			zap.Any("md", md),
		)
		return handler(ctx, req)
	}
}

// stream middlware for logging the request to stdout
func requestLoggingStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		panicForNilLogger()

		md, _ := metadata.FromIncomingContext(stream.Context())
		logger.Info("request from caller",
			zap.String("info.full_method", info.FullMethod),
			zap.Any("md", md),
		)

		return handler(srv, stream)
	}
}

// logs a response to stdout
func logResponse(_ context.Context, fullMethod string, response interface{}, err error) {
	panicForNilLogger()

	logger.Info("returned to caller",
		zap.Any("response", response),
		zap.String("info.full_method", fullMethod),
		zap.Error(err),
	)
}

// unary middleware for performing operations on the response
func responseUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		response, err := handler(ctx, req)
		defer func() {
			// things to do once there is a response
			logResponse(ctx, info.FullMethod, response, err)
			decrementActiveRequestCount()
			incrementRequestProcessed()
		}()

		return response, err
	}
}

// stream middleware for performing operation on the response
func responseStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		err := handler(srv, stream)
		defer func() {
			// things to do once there is a response
			logResponse(stream.Context(), info.FullMethod, "", err)
			// TODO :: should only happen after a stream is closed ...
			decrementActiveRequestCount()
			incrementRequestProcessed()
		}()

		return err
	}
}

// decrements the active_request_count on the health handler
func decrementActiveRequestCount() {
	svcHandlers.serviceHealth.DecrementActiveRequestCount()
}

// increments the request_processed count on the health handler
func incrementRequestProcessed() {
	svcHandlers.serviceHealth.IncrementRequestCount()
}

// increments the active_request_count on the health handler
func incrementActiveRequestCount() {
	svcHandlers.serviceHealth.IncrementActiveRequestCount()
}

// merges swagger json specs from multiple sources
func mergeSwaggerContent(from ...[]byte) []byte {
	sm := toucan_service.NewSwaggerMerger(
		toucan_service.WithFrom(from...),
		toucan_service.WithTags(
			[]spec.Tag{
				{
					TagProps: spec.TagProps{
						Description: "Health service API",
						Name:        "HealthService",
					},
				},
				// TODO :: Replace HelloWorld with your own
				{
					TagProps: spec.TagProps{
						Description: "Hello world service API",
						Name:        "HelloWorldService",
					},
				},
			},
		),
		toucan_service.WithInfo(&spec.Info{
			InfoProps: spec.InfoProps{
				Description: RootCmdLong,
				Title:       build.AppName,
				Version:     build.Version,
			},
		}),
	)

	mergedContent, warnings, err := sm.Merge()
	if err != nil {
		bail(err)
	}
	if len(warnings) != 0 {
		logger.Warn(fmt.Sprintf("swagger mixin warnings: %v", warnings))
	}

	return mergedContent
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// command flags
	serveCmd.Flags().BoolP(
		NoGRPCGatewayFlag,
		NoGRPCGatewayFlagShort,
		false,
		NoGRPCGatewayFlagUsage,
	)
	serveCmd.Flags().BoolP(
		NoSwaggerUIFlag,
		NoSwaggerUIFlagShort,
		false,
		NoSwaggerUIFlagUsage,
	)
	serveCmd.Flags().String(
		HostFlag,
		HostDefaultValue,
		HostFlagUsage,
	)
	serveCmd.Flags().Int32P(
		PortFlag,
		PortFlagShort,
		PortDefaultValue,
		PortFlagUsage,
	)

	err := gViper.BindPFlags(serveCmd.Flags())
	if err != nil {
		bail(err)
	}
}

func serveCmdFunc() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		defer server.Finished()

		ctx := context.Background()

		logger.Info(fmt.Sprintf("%s<version=%s, gitSha=%s>", build.AppName, build.Version, build.GitSha))

		logger.Info("go-kestrel-info", zap.String("details", cfg.AppInfo()))

		hnp := fmt.Sprintf("%s:%d", gViper.GetString(HostFlag), gViper.GetInt(PortFlag))

		// define the server options and middleware
		// TODO :: Replace HelloWorld with your own
		sopts := []server.Option{
			server.EnableReflection(),
			server.ProductionMiddleware(),
			// custom middleware
			server.AddUnaryServerMiddleware(requestLoggingUnaryInterceptor(), "request-unary-logger"),
			server.AddStreamServerMiddleware(requestLoggingStreamInterceptor(), "request-stream-logger"),
			server.AddUnaryServerMiddleware(requestCallBackUnaryInterceptor(
				func() { incrementActiveRequestCount() },
			), "counters-unary"),
			server.AddStreamServerMiddleware(requestCallBackStreamInterceptor(
				func() { incrementActiveRequestCount() },
			), "counters-stream"),
			server.AddUnaryServerMiddleware(responseUnaryInterceptor(), "response-unary-logger"),
			server.AddStreamServerMiddleware(responseStreamInterceptor(), "response-stream-logger"),
			server.RegisterGRPCGWEndpoint(health.RegisterHealthServiceHandlerFromEndpoint),
			server.RegisterGRPCGWEndpoint(apis.RegisterTreasureServiceHandlerFromEndpoint),
		}

		// enable swagger ui / specify the server option based on the cmd flags
		if !gViper.GetBool(NoSwaggerUIFlag) && !gViper.GetBool(NoGRPCGatewayFlag) {
			sopts = append(sopts, server.EnableSwagger(
				mergeSwaggerContent(
					health.SwaggerJSONContent,
					apis.SwaggerJSONContent,
				),
			))
		}

		// call new server with the options, etc
		s, err := server.NewServer(
			build.AppName,
			hnp,
			sopts...,
		)
		if err != nil {
			panic(fmt.Sprintf("cannot start %s; go-kestrel failure: %v", build.AppName, err))
		}

		// initialize the handlers
		// TODO :: Replace HelloWorld with your own
		svcHandlers = &serviceHandlers{
			serviceHealth: toucan_handlers.NewHealthServiceHandler(ctx, hnp, toucan_service.BuildInfo{
				AppName:   build.AppName,
				Version:   build.Version,
				GitSha:    build.GitSha,
				BuiltBy:   build.BuiltBy,
				BuildHost: build.BuildHost,
				BuildTime: build.BuildTime,
			}),
			serviceTreasure: handlers.NewTreasureServiceHandler(ctx),
		}

		// register the handlers with the gRPC server
		// TODO :: Replace HelloWorld with your own
		// TODO :: register each additional handlers
		health.RegisterHealthServiceServer(s.GRPCServer(), svcHandlers.serviceHealth)
		apis.RegisterTreasureServiceServer(s.GRPCServer(), svcHandlers.serviceTreasure)
		// start the server
		logger.Info(fmt.Sprintf("%s-done", build.AppName), zap.Error(s.Run()))
	}
}

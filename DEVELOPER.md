# piratetreasure

its cursed

A core test services project created from [`gpt_toucan`](https://github.com/netskope-qe/gpt_toucan).

## Building the source

Most of the developer environment needs for building this project have been packaged into docker images and incorporated into the build stucture. `make`, `drone`, and `docker` are the only tools required.

### Required tools

1. Install `docker`
2. Install the [`drone`](https://docs.drone.io/cli/install/) CLI into your `PATH`
3. Ensure your `make` binary is the GNU Make 4.x+

### Build steps

1. Build the protobuf with `make build-protobuf`.
2. Build the application binary with `make project`.

Notes:

* Build steps 1-2 are only required the first time building the project and when protobuf source changes happen.
The resulting generated code should be checked into source control.
* `make project` poduces a linux binary. If you are developing on a non-linux host, run `make build-piratetreasure` to produce a binary which is compatible with your host -- `go` is required.* Explore additional make targets with `make`.

## Project layout

```shell
api/proto            <- .proto files, generated {proto}.go, {proto}.gw.go, and swagger files.
build                <- build files
cmd/                 <- cobra and/or viper cli command defs
docs/arch/decisions  <- architectural decision docs related to the service
functests            <- funtional tests for the service
internal/handlers    <- service handlers
internal/build       <- service build info (dynamically updated during build)
```

## Starting the service

The service can be started via the CLI and the `serve` command:

```shell
./dist/piratetreasure serve --appname piratetreasure
```

For example:

```shell
$ ./dist/piratetreasure serve --appname piratetreasure
{"level":"info","ts":"2019-07-10T20:11:51.255-0700","caller":"tracing/tracing.go:50","msg":"tracedest","app":"piratetreasure","host":"localhost","port":"6831"}
{"level":"info","ts":"2019-07-10T20:11:51.256-0700","caller":"zap/logger.go:38","msg":"Initializing logging reporter\n","app":"piratetreasure"}
{"level":"info","ts":"2019-07-10T20:11:51.257-0700","caller":"cmd/serve.go:98","msg":"go-kestrel-info","app":"piratetreasure","details":"GoKestrel<version=v1.1.7>"}
...
```

## Accessing the service via Swagger UI

Open your favorite web browser to `http://{host}:{port}/swagger-ui/`

For this demo, that would be:

```shell
http://localhost:12345/swagger-ui/
```

You can disable the swagger UI by passing the `--no-swagger-ui` flag to the
`serve` command.

## Accessing the service via the CLI

The same CLI can be used to invoke the service operations:

```shell
./dist/piratetreasure invoke [command] [flags] [args]
```

For example:

```shell
$ ./dist/piratetreasure invoke health
Reply: {
  "serviceIdentity": {
    "appName": "piratetreasure",
    "appVersion": "v0.0.0",
    "builtBy": "doug",
    "gitSha": "be394fb687ab959f",
    "buildHost": "chopin",
    "buildTime": "0001-01-01T00:00:00Z"
  },
  "upTime": "43.906943752s",
  "listeningAddress": "0.0.0.0:12345"
}
```

## Accessing the service via curl

If the service is running with the gRPC gateway enabled, you can invoke the
operations in a REST compliant fashion. To do this with curl, for example:

```shell
$ curl -X GET "http://localhost:12345/v1/health" -H  "accept: application/json"
{"serviceIdentity":{"appName":"piratetreasure","appVersion":"v0.0.0","builtBy":"doug","gitSha":"be394fb687ab959f","buildHost":"chopin","buildTime":"0001-01-01T00:00:00Z"},"upTime":"121.753192152s","listeningAddress":"0.0.0.0:12345"}
```

## Accessing the service via grpcui

The open source tool `grpcui` allows you to invoke the gRPC service in a web UI similar to Swagger UI.
After `grpcui` is installed run the following;

```shell
grpcui -plaintext {host}:{port}
```

Then open the link which is displayed in your favorite web browser. For example;

```shell
$ grpcui -plaintext localhost:12345
gRPC Web UI available at http://127.0.0.1:33445/
```

## Accessing the service via grpcurl

The open source tool `grpcurl` allows you to invoke the gRPC service in a fashion similar to curl. After `grpcurl` is installed run the following;

```shell
grpcurl -plaintext {host}:{port} [command] [symbol]
```

For example:

```shell
$ grpcurl -plaintext localhost:12345 HealthService.Health
{
  "serviceIdentity": {
    "appName": "piratetreasure",
    "appVersion": "v0.0.0",
    "builtBy": "doug",
    "gitSha": "be394fb687ab959f",
    "buildHost": "chopin",
    "buildTime": "0001-01-01T00:00:00Z"
  },
  "upTime": "43.906943752s",
  "listeningAddress": "0.0.0.0:12345"
}
```

## Configuring the service

The `serve` command has the following options;

```shell
$ ./dist/piratetreasure serve -h
Starts the service on the specified {host}:{port} and listens for requests

Usage:
  piratetreasure serve [flags]

Flags:
  -h, --help              help for serve
      --host string       the host to listen on (default "0.0.0.0")
  -g, --no-grpc-gateway   disables the REST<->gRPC gateway and the Swagger UI
  -s, --no-swagger-ui     disables the Swagger UI (http(s)://{host}:{port}/swagger-ui/)
  -p, --port int32        the port to bind to (default 12345)

Global Flags:
      --appname string     Name of application to (used for observability) (default "unknown")
      --config string      config file (default is $HOME/.piratetreasure.yaml)
      --loglvl string      Log level to use for logging (default "info")
      --tracedest string   Destination host/ip to send tracing data (default "localhost:6831")
      --tracetags string   Comma-separated list of name=value tags to send with traces
```

## Metrics

You can start the metrics infrastructure with `make start-metrics`.

Then run some traffic through the piratetreasure and:

* For Prometheus: [http://localhost:9090](http://localhost:9090)
* For Jaeger tracing: [http://localhost:16686](http://localhost:16686)

## Generate some load

You can use the third party application [`ghz`](https://ghz.sh/) to generate some load.

For example;

```shell
ghz --insecure --call HealthService.Health localhost:12345
```

/*
Copyright © 2019-2021 Netskope
*/

package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/netskope/go-kestrel/pkg/log"

	"github.com/netskope/piratetreasure/api/proto/piratetreasure"
)

var (
	swaggerLogger = log.NewLogger("SwaggerJSONHandler")
)

// SwaggerHTTPHandler returns a handler that serves the generated swagger.json for this service.
// This function just wraps the file from bindata into an handler function.
func SwaggerHTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		_, err := io.Copy(w, bytes.NewReader(piratetreasure.SwaggerJSONContent))
		if err != nil {
			swaggerLogger.Error(fmt.Sprintf("error occurred loading swaggerjson content, %s", err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

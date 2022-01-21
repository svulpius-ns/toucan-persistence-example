/*
Copyright © 2019-2020 Netskope
*/

package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/netskope/go-kestrel/pkg/test"

	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	buf          strings.Builder
	forceFailure bool
	rc           int
}

func (m *mockWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockWriter) Write(b []byte) (int, error) {
	if m.forceFailure {
		return 0, fmt.Errorf("mockWriter.Write failure")
	}
	return m.buf.Write(b)
}

func (m *mockWriter) WriteHeader(statusCode int) {
	m.rc = statusCode
}

func TestSwaggerHTTPHandler(t *testing.T) {
	test.InitKestrelForTest()
	w := mockWriter{}
	handler := SwaggerHTTPHandler()
	handler(&w, &http.Request{})

	assert.True(t, strings.Contains(w.buf.String(), "swagger"))
}

func TestSwaggerHTTPHandler_Failure(t *testing.T) {
	test.InitKestrelForTest()
	w := mockWriter{forceFailure: true}
	handler := SwaggerHTTPHandler()
	handler(&w, &http.Request{})

	assert.Equal(t, "", w.buf.String())
	assert.Equal(t, http.StatusInternalServerError, w.rc)
}

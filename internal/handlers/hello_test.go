/*
Copyright © 2019-2020 Netskope
*/

package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/netskope/go-kestrel/pkg/test"
)

func TestNewHelloWorldServiceHandler(t *testing.T) {
	test.InitKestrelForTest()

	sctx := context.Background()
	h := NewHelloWorldServiceHandler(sctx)
	assert.NotNil(t, h)
}

func TestSayHello(t *testing.T) {
	test.InitKestrelForTest()

	sctx := context.Background()
	rctx := context.Background()
	hw := NewHelloWorldServiceHandler(sctx)
	helloResponse, err := hw.SayHello(rctx, &empty.Empty{})
	assert.Nil(t, err)
	assert.NotEmpty(t, helloResponse)
	assert.Equal(t, "Hello World!", helloResponse.GetHello())
}

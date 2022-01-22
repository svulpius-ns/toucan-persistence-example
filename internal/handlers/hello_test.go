/*
Copyright © 2019-2020 Netskope
*/

package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/netskope/go-kestrel/pkg/test"
	apis "github.com/netskope/piratetreasure/api/proto/piratetreasure"
)

func TestNewHelloWorldServiceHandler(t *testing.T) {
	test.InitKestrelForTest()

	sctx := context.Background()
	h := NewTreasureServiceHandler(sctx)
	assert.NotNil(t, h)
}

func TestSayHello(t *testing.T) {
	test.InitKestrelForTest()

	sctx := context.Background()
	rctx := context.Background()
	hw := NewTreasureServiceHandler(sctx)
	helloResponse, err := hw.ListTreasure(rctx, &apis.ListTreasureRequest{})
	assert.Nil(t, err)
	assert.NotEmpty(t, helloResponse)
	assert.Equal(t, "Hello World!", helloResponse.Treasure)
}

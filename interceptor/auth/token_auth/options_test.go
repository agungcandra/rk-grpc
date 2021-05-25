// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkgrpctokenauth

import (
	"github.com/rookie-ninja/rk-grpc/interceptor/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithEntryNameAndType_HappyCase(t *testing.T) {
	set := newOptionSet(rkgrpcctx.RpcTypeUnaryServer,
		WithEntryNameAndType("ut-entry-name", "ut-entry"))

	assert.Equal(t, "ut-entry-name", set.EntryName)
	assert.Equal(t, "ut-entry", set.EntryType)
	assert.Equal(t, set,
		optionsMap[rkgrpcctx.ToOptionsKey("ut-entry-name", rkgrpcctx.RpcTypeUnaryServer)])
}

func TestWithToken_HappyCase(t *testing.T) {
	token := "token"

	set := newOptionSet(rkgrpcctx.RpcTypeUnaryServer,
		WithToken(token, true))

	assert.True(t, set.tokens["token"])
}

func TestOptionSet_Authorized_HappyCase(t *testing.T) {
	set := newOptionSet(rkgrpcctx.RpcTypeUnaryServer,
		WithToken("token-1", true),
		WithToken("token-2", false))

	assert.False(t, set.Authorized("token-1"))
	assert.True(t, set.Authorized("token-2"))
}
// Copyright (c) 2018 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package apikey

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
	rsapkg "github.com/northwesternmutual/kanali/pkg/rsa"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestApiKeyAdd(t *testing.T) {
	lvl := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	core, logs := observer.New(lvl)
	defer log.SetLogger(zap.New(core)).Restore()

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	ctlr := NewController(priv)
	encryptedKey, _ := rsapkg.Encrypt([]byte("abc123"), &priv.PublicKey, rsapkg.Base64Encode(), rsapkg.WithEncryptionLabel(rsapkg.EncryptionLabel))
	apikey := builder.NewApiKey("foo").WithRevision(v2.RevisionStatusActive, encryptedKey).NewOrDie()

	assert.True(t, store.ApiKeyStore().IsEmpty())
	ctlr.OnAdd(apikey)
	assert.Equal(t, 1, logs.FilterMessageSnippet("added").Len())
	assert.NotNil(t, store.ApiKeyStore().Get("abc123"))

	ctlr.OnAdd(builder.NewApiKey("bar").WithRevision(v2.RevisionStatusActive, []byte("foo")).NewOrDie())
	assert.Equal(t, 1, logs.FilterMessageSnippet("illegal").Len())

	ctlr.OnAdd(nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())

	ctlr.OnAdd("foo")
	assert.Equal(t, 2, logs.FilterMessageSnippet("malformed").Len())
}

func TestApiKeyUpdate(t *testing.T) {

}

func TestApiKeyDelete(t *testing.T) {

}

func TestDecryptApiKey(t *testing.T) {

}

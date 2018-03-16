// Copyright (c) 2017 Northwestern Mutual.
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
	"crypto/rsa"
	"errors"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"

	rsautils "github.com/northwesternmutual/kanali/pkg/rsa"
)

type apiKeyController struct {
	decryptionKey *rsa.PrivateKey
}

func NewController(decryptionKey *rsa.PrivateKey) cache.ResourceEventHandler {
	return &apiKeyController{
		decryptionKey: decryptionKey,
	}
}

func (ctlr *apiKeyController) OnAdd(obj interface{}) {
	logger := log.WithContext(nil)
	key, ok := obj.(*v2.ApiKey)
	if !ok || key == nil {
		logger.Error("received malformed ApiKey from k8s apiserver")
		return
	}
	keyClone, err := ctlr.decryptApiKey(key)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	store.ApiKeyStore().Set(keyClone)
	logger.With(
		zap.String(tags.KanaliApiKeyName, keyClone.GetName()),
	).Debug("added ApiKey")
}

func (ctlr *apiKeyController) OnUpdate(old interface{}, new interface{}) {
	logger := log.WithContext(nil)
	newKey, ok := new.(*v2.ApiKey)
	if !ok {
		logger.Error("received malformed ApiKey from k8s apiserver")
		return
	}
	oldKey, ok := old.(*v2.ApiKey)
	if !ok {
		logger.Error("received malformed ApiKey from k8s apiserver")
		return
	}
	newKeyClone, err := ctlr.decryptApiKey(newKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	oldKeyClone, err := ctlr.decryptApiKey(oldKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	store.ApiKeyStore().Update(oldKeyClone, newKeyClone)
	logger.With(
		zap.String(tags.KanaliApiKeyName, newKeyClone.GetName()),
	).Debug("updated ApiKey")
}

func (ctlr *apiKeyController) OnDelete(obj interface{}) {
	logger := log.WithContext(nil)
	key, ok := obj.(*v2.ApiKey)
	if !ok {
		logger.Error("received malformed ApiKey from k8s apiserver")
		return
	}
	key, err := ctlr.decryptApiKey(key)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if result := store.ApiKeyStore().Delete(key); result != nil {
		logger.With(
			zap.String(tags.KanaliApiKeyName, result.GetName()),
		).Debug("deleted ApiKey")
	}
}

func (ctlr *apiKeyController) decryptApiKey(key *v2.ApiKey) (*v2.ApiKey, error) {
	if ctlr.decryptionKey == nil {
		return nil, errors.New("decryption key not present")
	}

	clone := key.DeepCopy()

	for i, revision := range clone.Spec.Revisions {
		unencryptedApiKey, err := rsautils.Decrypt([]byte(revision.Data), ctlr.decryptionKey,
			rsautils.Base64Decode(),
			rsautils.WithEncryptionLabel(rsautils.EncryptionLabel),
		)
		if err != nil {
			return nil, err
		}
		clone.Spec.Revisions[i].Data = string(unencryptedApiKey)
	}
	return clone, nil
}

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
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	informers "github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"k8s.io/client-go/tools/cache"
)

type ApiKeyController struct {
	apikeys       informers.ApiKeyInformer
	decryptionKey *rsa.PrivateKey
}

func NewApiKeyController(apikeys informers.ApiKeyInformer, decryptionKey *rsa.PrivateKey) *ApiKeyController {

	ctlr := &ApiKeyController{
		decryptionKey: decryptionKey,
	}

	ctlr.apikeys = apikeys
	apikeys.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    ctlr.apiKeyAdd,
			UpdateFunc: ctlr.apiKeyUpdate,
			DeleteFunc: ctlr.apiKeyDelete,
		},
		5*time.Minute,
	)

	return ctlr
}

func (ctlr *ApiKeyController) Run(stopCh <-chan struct{}) {
	ctlr.apikeys.Informer().Run(stopCh)
}

func (ctlr *ApiKeyController) apiKeyAdd(obj interface{}) {
	logger := logging.WithContext(nil)
	key, ok := obj.(*v2.ApiKey)
	if !ok {
		logger.Error("received malformed ApiKey from k8s apiserver")
		return
	}
	if err := ctlr.decryptApiKey(key); err != nil {
		logger.Error(err.Error())
		return
	}
	store.ApiKeyStore.Set(*key)
	logger.Debug(fmt.Sprintf("added ApiKey %s", key.ObjectMeta.Name))
}

func (ctlr *ApiKeyController) apiKeyUpdate(old interface{}, new interface{}) {
	logger := logging.WithContext(nil)
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
	if err := ctlr.decryptApiKey(newKey); err != nil {
		logger.Error(err.Error())
		return
	}
	if err := ctlr.decryptApiKey(oldKey); err != nil {
		logger.Error(err.Error())
		return
	}
	store.ApiKeyStore.Update(*oldKey, *newKey)
	logger.Debug(fmt.Sprintf("updated ApiKey %s", newKey.ObjectMeta.Name))
}

func (ctlr *ApiKeyController) apiKeyDelete(obj interface{}) {
	logger := logging.WithContext(nil)
	key, ok := obj.(*v2.ApiKey)
	if !ok {
		logger.Error("received malformed ApiKey from k8s apiserver")
		return
	}
	if err := ctlr.decryptApiKey(key); err != nil {
		logger.Error(err.Error())
		return
	}
	result, _ := store.ApiKeyStore.Delete(*key)
	if result != nil {
		result := result.(v2.ApiKey)
		logger.Debug(fmt.Sprintf("deleted ApiKey %s", result.ObjectMeta.Name))
	}
}

func (ctlr *ApiKeyController) decryptApiKey(key *v2.ApiKey) error {
	if ctlr.decryptionKey == nil {
		return errors.New("decryption key not present")
	}

	for _, revision := range key.Spec.Revisions {
		cipherText, err := hex.DecodeString(revision.Data)
		if err != nil {
			return err
		}
		unencryptedApiKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, ctlr.decryptionKey, cipherText, []byte("kanali"))
		if err != nil {
			return err
		}
		revision.Data = string(unencryptedApiKey)
	}
	return nil
}

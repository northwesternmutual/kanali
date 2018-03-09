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
	"context"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
	informers "github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	rsautils "github.com/northwesternmutual/kanali/pkg/rsa"
)

type ApiKeyController struct {
	apikeys       informers.ApiKeyInformer
	decryptionKey *rsa.PrivateKey
	clientset     *versioned.Clientset
}

func NewApiKeyController(apikeys informers.ApiKeyInformer, clientset *versioned.Clientset, decryptionKey *rsa.PrivateKey) *ApiKeyController {

	ctlr := &ApiKeyController{
		decryptionKey: decryptionKey,
		clientset:     clientset,
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

func (ctlr *ApiKeyController) Run(ctx context.Context) error {
	ctlr.apikeys.Informer().Run(ctx.Done())
	return nil
}

func (ctlr *ApiKeyController) Close(error) error {
	return nil
}

func (ctlr *ApiKeyController) apiKeyAdd(obj interface{}) {
	logger := log.WithContext(nil)
	key, ok := obj.(*v2.ApiKey)
	if !ok || key == nil {
		logger.Error("received malformed ApiKey from k8s apiserver")
		return
	}
	keyClone, err := ctlr.decryptApiKey(key)
	if err != nil {
		if err := ctlr.clientset.KanaliV2().ApiKeys().Delete(key.GetName(), &v1.DeleteOptions{}); err != nil {
			logger.Error(err.Error())
		}
		logger.Error(err.Error())
		return
	}
	store.ApiKeyStore().Set(keyClone)
	logger.With(
		zap.String(tags.KanaliApiKeyName, keyClone.GetName()),
	).Debug("added ApiKey")
}

func (ctlr *ApiKeyController) apiKeyUpdate(old interface{}, new interface{}) {
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
		if err := ctlr.clientset.KanaliV2().ApiKeys().Delete(newKeyClone.GetName(), &v1.DeleteOptions{}); err != nil {
			logger.Error(err.Error())
		}
		logger.Error(err.Error())
		return
	}
	oldKeyClone, err := ctlr.decryptApiKey(oldKey)
	if err != nil {
		if err := ctlr.clientset.KanaliV2().ApiKeys().Delete(oldKeyClone.GetName(), &v1.DeleteOptions{}); err != nil {
			logger.Error(err.Error())
		}
		logger.Error(err.Error())
		return
	}
	store.ApiKeyStore().Update(oldKeyClone, newKeyClone)
	logger.With(
		zap.String(tags.KanaliApiKeyName, newKeyClone.GetName()),
	).Debug("updated ApiKey")
}

func (ctlr *ApiKeyController) apiKeyDelete(obj interface{}) {
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

func (ctlr *ApiKeyController) decryptApiKey(key *v2.ApiKey) (*v2.ApiKey, error) {
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

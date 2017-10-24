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
	"github.com/northwesternmutual/kanali/pkg/store"
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

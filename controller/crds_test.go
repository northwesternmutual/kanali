package controller

import (
  "testing"

	"github.com/stretchr/testify/assert"
  "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
)

func TestCreateTPRs(t *testing.T) {
	ctlr := Controller{
		RESTClient: nil,
		ClientSet:  nil,
    APIExtensionsV1beta1Interface: fake.NewSimpleClientset().ApiextensionsV1beta1(),
	}

  err := ctlr.CreateCRDs()
  assert.Nil(t, err)

  proxyCRD, err := ctlr.APIExtensionsV1beta1Interface.CustomResourceDefinitions().Get("apiproxies.kanali.io", v1.GetOptions{})
  assert.Nil(t, err)
  assert.Equal(t, proxyCRD.ObjectMeta.Name, "apiproxies.kanali.io")

  bindingCRD, err := ctlr.APIExtensionsV1beta1Interface.CustomResourceDefinitions().Get("apikeybindings.kanali.io", v1.GetOptions{})
  assert.Nil(t, err)
  assert.Equal(t, bindingCRD.ObjectMeta.Name, "apikeybindings.kanali.io")

  keyCRD, err := ctlr.APIExtensionsV1beta1Interface.CustomResourceDefinitions().Get("apikeys.kanali.io", v1.GetOptions{})
  assert.Nil(t, err)
  assert.Equal(t, keyCRD.ObjectMeta.Name, "apikeys.kanali.io")
}
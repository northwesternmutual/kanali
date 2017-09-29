package controller

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

const (
  KanaliGroupName = "kanali.io"
)

var apiProxyCRD = &apiextensionsv1beta1.CustomResourceDefinition{
  ObjectMeta: metav1.ObjectMeta{
    Name: fmt.Sprintf("apiproxies.%s", KanaliGroupName),
  },
  Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
    Group:   KanaliGroupName,
    Version: "v1",
    Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
      Plural: "apiproxies",
      Singular: "apiproxy",
      ShortNames: []string{
        "ap",
        "proxies",
      },
      Kind: "ApiProxy",
      ListKind: "ApiProxyList",
    },
    Scope: apiextensionsv1beta1.NamespaceScoped,
  },
}

var apiKeyBindingCRD = &apiextensionsv1beta1.CustomResourceDefinition{
  ObjectMeta: metav1.ObjectMeta{
    Name: fmt.Sprintf("apikeybindings.%s", KanaliGroupName),
  },
  Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
    Group:   KanaliGroupName,
    Version: "v1",
    Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
      Plural: "apikeybindings",
      Singular: "apikeybinding",
      ShortNames: []string{
        "akb",
        "bindings",
      },
      Kind: "ApiKeyBinding",
      ListKind: "ApiKeyBindingList",
    },
    Scope: apiextensionsv1beta1.NamespaceScoped,
  },
}

var apiKeyCRD = &apiextensionsv1beta1.CustomResourceDefinition{
  ObjectMeta: metav1.ObjectMeta{
    Name: fmt.Sprintf("apikeys.%s", KanaliGroupName),
  },
  Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
    Group:   KanaliGroupName,
    Version: "v1",
    Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
      Plural: "apikeys",
      Singular: "apikey",
      ShortNames: []string{
        "ak",
        "keys",
      },
      Kind: "ApiKey",
      ListKind: "ApiKeyList",
    },
    Scope: apiextensionsv1beta1.ClusterScoped,
  },
}

// CreateCRDs will create all the CRDs that Kanali requires
func (c *Controller) CreateCRDs() error {
  return doCreateCRDs(c, apiProxyCRD, apiKeyBindingCRD, apiKeyCRD)
}

func doCreateCRDs(c *Controller, crds ...*apiextensionsv1beta1.CustomResourceDefinition) error {
  for _, crd := range crds {
    _, err := c.APIExtensionsV1beta1Interface.CustomResourceDefinitions().Create(crd)
    if err != nil && !errors.IsAlreadyExists(err) {
      return fmt.Errorf("Failed to create CRD: %v", err)
    }
  }

  return nil
}

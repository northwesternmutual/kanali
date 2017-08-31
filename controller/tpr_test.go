package controller

import (
	"errors"
	"net/http"
	"testing"

  "k8s.io/kubernetes/pkg/api"
	"github.com/stretchr/testify/assert"
	e "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
  "k8s.io/kubernetes/pkg/apis/extensions"
  "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/fake"
)

func TestCreateTPRs(t *testing.T) {
  ctlr := Controller{
    RestClient: nil,
    ClientSet: fake.NewSimpleClientset(),
    MasterHost: "foo.bar.com",
  }

  err := ctlr.CreateTPRs()
  assert.Nil(t, err)
  err = ctlr.CreateTPRs()
  assert.Nil(t, err)
  resource, _ := ctlr.ClientSet.Extensions().ThirdPartyResources().Get("api-proxy.kanali.io")
  assert.Equal(t, resource, &extensions.ThirdPartyResource{
    ObjectMeta: api.ObjectMeta{
      Name: "api-proxy.kanali.io",
    },
    Versions: []extensions.APIVersion{
      {
        Name: "v1",
      },
    },
    Description: "api proxy TPR",
  })
  resource, _ = ctlr.ClientSet.Extensions().ThirdPartyResources().Get("api-key.kanali.io")
  assert.Equal(t, resource, &extensions.ThirdPartyResource{
    ObjectMeta: api.ObjectMeta{
      Name: "api-key.kanali.io",
    },
    Versions: []extensions.APIVersion{
      {
        Name: "v1",
      },
    },
    Description: "api key TPR",
  })
  resource, _ = ctlr.ClientSet.Extensions().ThirdPartyResources().Get("api-key-binding.kanali.io")
  assert.Equal(t, resource, &extensions.ThirdPartyResource{
    ObjectMeta: api.ObjectMeta{
      Name: "api-key-binding.kanali.io",
    },
    Versions: []extensions.APIVersion{
      {
        Name: "v1",
      },
    },
    Description: "api key binding TPR",
  })
}

func TestIsKubernetesResourceAlreadyExistError(t *testing.T) {
	assert.False(t, isKubernetesResourceAlreadyExistError(errors.New("test error")))
	se := e.StatusError{
		ErrStatus: unversioned.Status{
			Code:   http.StatusConflict,
			Reason: unversioned.StatusReasonAlreadyExists,
		},
	}
	assert.True(t, isKubernetesResourceAlreadyExistError(&se))
	se.ErrStatus.Code = http.StatusNotFound
	assert.False(t, isKubernetesResourceAlreadyExistError(&se))
}
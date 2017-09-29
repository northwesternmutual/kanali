package controller

import (
	"testing"

	"github.com/northwesternmutual/kanali/crds"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestAddKnownTypes(t *testing.T) {
	scheme := runtime.NewScheme()
	err := addKnownTypes(scheme)
	assert.Nil(t, err)

	types := scheme.AllKnownTypes()
	_, ok := types[schema.GroupVersionKind{
		Group:   crds.KanaliGroupName,
		Version: "v1",
		Kind:    "ApiProxy",
	}]
	assert.True(t, ok)

	_, ok = types[schema.GroupVersionKind{
		Group:   crds.KanaliGroupName,
		Version: "v1",
		Kind:    "ApiKey",
	}]
	assert.True(t, ok)

	_, ok = types[schema.GroupVersionKind{
		Group:   crds.KanaliGroupName,
		Version: "v1",
		Kind:    "ApiKeyBinding",
	}]
	assert.True(t, ok)
}

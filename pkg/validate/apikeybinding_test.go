package validate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestIsValidApiKeyBinding(t *testing.T) {
	assert.Error(t, (&validation{}).IsValidApiKeyBinding(nil))

	apikeybinding, _ := json.Marshal(builder.NewApiKeyBinding("foo", "bar").NewOrDie())
	assert.Nil(t, (&validation{}).IsValidApiKeyBinding(apikeybinding))
}

func TestIsValidApiKeyBindingList(t *testing.T) {
	assert.Error(t, (&validation{}).IsValidApiKeyBindingList(nil))

	list, _ := json.Marshal(&v2.ApiKeyBindingList{
		Items: []v2.ApiKeyBinding{
			*builder.NewApiKeyBinding("foo", "bar").NewOrDie(),
		},
	})
	assert.Nil(t, (&validation{}).IsValidApiKeyBindingList(list))
}

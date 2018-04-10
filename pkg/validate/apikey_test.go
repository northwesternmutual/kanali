package validate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestIsValidApiKey(t *testing.T) {
	assert.Error(t, (&validation{}).IsValidApiKey(nil))

	apikey, _ := json.Marshal(builder.NewApiKey("foo").NewOrDie())
	assert.Nil(t, (&validation{}).IsValidApiKey(apikey))
}

func TestIsValidApiKeyList(t *testing.T) {
	assert.Error(t, (&validation{}).IsValidApiKeyList(nil))

	list, _ := json.Marshal(&v2.ApiKeyList{
		Items: []v2.ApiKey{
			*builder.NewApiKey("foo").NewOrDie(),
		},
	})
	assert.Nil(t, (&validation{}).IsValidApiKeyList(list))
}

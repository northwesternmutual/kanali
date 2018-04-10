package validate

import (
	"encoding/json"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

func (v *validation) IsValidApiKey(data []byte) error {
	apikey := new(v2.ApiKey)
	if err := json.Unmarshal(data, apikey); err != nil {
		return err
	}
	return v.isValidApiKey(apikey)
}

func (v *validation) isValidApiKey(apikey *v2.ApiKey) error {
	// no dynamic validation needed at this time
	return nil
}

func (v *validation) IsValidApiKeyList(data []byte) error {
	list := new(v2.ApiKeyList)
	if err := json.Unmarshal(data, list); err != nil {
		return err
	}
	return v.isValidApiKeyList(list)
}

func (v *validation) isValidApiKeyList(list *v2.ApiKeyList) error {
	for _, apikey := range list.Items {
		if err := v.isValidApiKey(&apikey); err != nil {
			return err
		}
	}
	return nil
}

package validate

import (
	"encoding/json"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

func (v *validation) IsValidApiKeyBinding(data []byte) error {
	apikeybinding := new(v2.ApiKeyBinding)
	if err := json.Unmarshal(data, apikeybinding); err != nil {
		return err
	}
	return v.isValidApiKeyBinding(apikeybinding)
}

func (v *validation) isValidApiKeyBinding(apikeybinding *v2.ApiKeyBinding) error {
	// no dynamic validation needed at this time
	return nil
}

func (v *validation) IsValidApiKeyBindingList(data []byte) error {
	list := new(v2.ApiKeyBindingList)
	if err := json.Unmarshal(data, list); err != nil {
		return err
	}
	return v.isValidApiKeyBindingList(list)
}

func (v *validation) isValidApiKeyBindingList(list *v2.ApiKeyBindingList) error {
	for _, apikeybinding := range list.Items {
		if err := v.isValidApiKeyBinding(&apikeybinding); err != nil {
			return err
		}
	}
	return nil
}

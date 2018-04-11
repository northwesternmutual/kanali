// Copyright (c) 2018 Northwestern Mutual.
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

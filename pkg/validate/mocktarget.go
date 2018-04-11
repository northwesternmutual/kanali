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

func (v *validation) IsValidMockTarget(data []byte) error {
	mocktarget := new(v2.MockTarget)
	if err := json.Unmarshal(data, mocktarget); err != nil {
		return err
	}
	return v.isValidMockTarget(mocktarget)
}

func (v *validation) isValidMockTarget(mocktarget *v2.MockTarget) error {
	// no dynamic validation needed at this time
	return nil
}

func (v *validation) IsValidMockTargetList(data []byte) error {
	list := new(v2.MockTargetList)
	if err := json.Unmarshal(data, list); err != nil {
		return err
	}
	return v.isValidMockTargetList(list)
}

func (v *validation) isValidMockTargetList(list *v2.MockTargetList) error {
	for _, mocktarget := range list.Items {
		if err := v.isValidMockTarget(&mocktarget); err != nil {
			return err
		}
	}
	return nil
}

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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestIsValidMockTarget(t *testing.T) {
	assert.Error(t, (&validation{}).IsValidMockTarget(nil))

	mocktarget, _ := json.Marshal(builder.NewMockTarget("foo", "bar").NewOrDie())
	assert.Nil(t, (&validation{}).IsValidMockTarget(mocktarget))
}

func TestIsValidMockTargetList(t *testing.T) {
	assert.Error(t, (&validation{}).IsValidMockTargetList(nil))

	list, _ := json.Marshal(&v2.MockTargetList{
		Items: []v2.MockTarget{
			*builder.NewMockTarget("foo", "bar").NewOrDie(),
		},
	})
	assert.Nil(t, (&validation{}).IsValidMockTargetList(list))
}

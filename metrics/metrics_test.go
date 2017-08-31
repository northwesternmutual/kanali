// Copyright (c) 2017 Northwestern Mutual.
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

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	f := &Metrics{}
	f.Add(Metric{"nameOne", "valueOne", false})
	f.Add(Metric{"nameTwo", "valueTwo", true}, Metric{"nameThree", "valueThree", false})
	assert.Equal(t, len(*f), 3)
	assert.False(t, (*f)[0].Index)
	assert.True(t, (*f)[1].Index)
	assert.False(t, (*f)[2].Index)
	assert.Equal(t, (*f)[0].Name, "nameOne")
	assert.Equal(t, (*f)[1].Name, "nameTwo")
	assert.Equal(t, (*f)[2].Name, "nameThree")
	assert.Equal(t, (*f)[0].Value, "valueOne")
	assert.Equal(t, (*f)[1].Value, "valueTwo")
	assert.Equal(t, (*f)[2].Value, "valueThree")
}

func TestGet(t *testing.T) {
	f := &Metrics{}
	f.Add(Metric{"nameOne", "valueOne", false})
	f.Add(Metric{"nameTwo", "valueTwo", true}, Metric{"nameThree", "valueThree", false})
	assert.Nil(t, f.Get("nameFour"))
	assert.Equal(t, f.Get("nameOne"), &Metric{"nameOne", "valueOne", false})
	assert.Equal(t, f.Get("nameTwo"), &Metric{"nameTwo", "valueTwo", true})
	assert.Equal(t, f.Get("nameThree"), &Metric{"nameThree", "valueThree", false})
}

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

package tracer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoParseConfig(t *testing.T) {
	config := []byte(`
---
sampler:
  type: const
  samplingServerURL: foo.bar.com
  param: 1
reporter:
  logSpans: true
  localAgentHostPort: "foo.bar.com:5775"
  `)

	cfg, err := doParseConfig(bytes.NewReader(config))
	assert.Nil(t, err)
	assert.Equal(t, "const", cfg.Sampler.Type)
	assert.Equal(t, float64(1), cfg.Sampler.Param)
	assert.Equal(t, "foo.bar.com", cfg.Sampler.SamplingServerURL)
	assert.Equal(t, true, cfg.Reporter.LogSpans)
	assert.Equal(t, "foo.bar.com:5775", cfg.Reporter.LocalAgentHostPort)

	_, err = doParseConfig(nil)
	assert.NotNil(t, err)

	_, err = doParseConfig(bytes.NewReader([]byte("!@#$%^&*()")))
	assert.NotNil(t, err)
}

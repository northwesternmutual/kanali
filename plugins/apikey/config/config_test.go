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

package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	generatedConfig, err := New(nil)
	if err == nil {
		t.Fail()
	}
	if err.Error() != "found nil config" {
		t.Fail()
	}
	if generatedConfig != nil {
		t.Fail()
	}

	_, err = New(map[string]string{
		"apiKeyBindingName": "apikey",
	})
	if err != nil {
		t.Fail()
	}
}

func TestBuildConfig(t *testing.T) {
	goodSourceConfig := map[string]string{
		"apiKeyBindingName": "apikey",
	}
	badSourceConfig := map[string]string{
		"foo": "bar",
	}

	generatedConfig, err := buildConfig(goodSourceConfig)
	if err != nil {
		t.Fail()
	}
	if generatedConfig == nil {
		t.Fail()
	}
	if generatedConfig.ApiKeyBindingName != "apikey" {
		t.Fail()
	}

	generatedConfig, err = buildConfig(badSourceConfig)
	if err == nil {
		t.Fail()
	}
	if err.Error() != ("required field apiKeyBindingName not found in config") {
		t.Fail()
	}
	if generatedConfig != nil {
		t.Fail()
	}
}

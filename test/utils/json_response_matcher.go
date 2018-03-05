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

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/onsi/gomega/types"
)

type representJSONMatcher struct {
	expected interface{}
}

func RepresentJSONifiedObject(expected interface{}) types.GomegaMatcher {
	return &representJSONMatcher{
		expected: expected,
	}
}

func (matcher *representJSONMatcher) Match(actual interface{}) (success bool, err error) {
	response, ok := actual.(*http.Response)
	if !ok {
		return false, fmt.Errorf("RepresentJSONifiedObject matcher expects an http.Response")
	}

	pointerToObjectOfExpectedType := reflect.New(reflect.TypeOf(matcher.expected)).Interface()
	err = json.NewDecoder(response.Body).Decode(pointerToObjectOfExpectedType)

	if err != nil {
		return false, fmt.Errorf("Failed to decode JSON: %s", err.Error())
	}

	decodedObject := reflect.ValueOf(pointerToObjectOfExpectedType).Elem().Interface()

	return reflect.DeepEqual(decodedObject, matcher.expected), nil
}

func (matcher *representJSONMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain the JSON representation of\n\t%#v", actual, matcher.expected)
}

func (matcher *representJSONMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to contain the JSON representation of\n\t%#v", actual, matcher.expected)
}

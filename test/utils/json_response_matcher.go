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

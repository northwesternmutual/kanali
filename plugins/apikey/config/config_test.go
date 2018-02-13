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

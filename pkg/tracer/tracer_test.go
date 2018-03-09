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

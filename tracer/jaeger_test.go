package tracer

import (
	"bytes"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCustomLoggerError(t *testing.T) {
	logger := customLogger{}

	writerOne := new(bytes.Buffer)
	writerTwo := new(bytes.Buffer)
	logrus.SetOutput(writerOne)
	logger.Error("custom error message")
	logrus.SetOutput(writerTwo)
	logrus.Error("custom error message")
	assert.Equal(t, writerOne.String(), writerTwo.String())

	writerOne = new(bytes.Buffer)
	writerTwo = new(bytes.Buffer)
	logrus.SetOutput(writerOne)
	logger.Infof("custom %s message", "info")
	logrus.SetOutput(writerTwo)
	logrus.Infof("custom %s message", "info")
	assert.Equal(t, writerOne.String(), writerTwo.String())
}

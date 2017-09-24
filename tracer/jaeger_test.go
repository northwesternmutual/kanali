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

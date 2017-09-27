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

package server

import (
	"testing"

	"github.com/northwesternmutual/kanali/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetKanaliPort(t *testing.T) {
	assert.Equal(t, getKanaliPort(), 80)

	viper.Set(config.FlagServerPort.GetLong(), 12345)
	assert.Equal(t, getKanaliPort(), 12345)

	viper.Set(config.FlagServerPort.GetLong(), 0)
	viper.Set(config.FlagTLSCertFile.GetLong(), "hi")
	viper.Set(config.FlagTLSKeyFile.GetLong(), "bye")
	assert.Equal(t, getKanaliPort(), 443)
}

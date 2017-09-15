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

package cmd

import (
	"testing"

	"github.com/northwesternmutual/kanali/config"
	"github.com/stretchr/testify/assert"
)

func TestStartCmdInit(t *testing.T) {
	assert.Equal(t, len(RootCmd.Commands()), 2)
	assert.Equal(t, RootCmd.Commands()[0], startCmd)

	for _, f := range *(config.Flags) {
		assert.Equal(t, startCmd.Flag(f.GetLong()).Name, f.GetLong())
		assert.Equal(t, startCmd.Flag(f.GetLong()).Shorthand, f.GetShort())
		assert.Equal(t, startCmd.Flag(f.GetLong()).Usage, f.GetUsage())
	}
}

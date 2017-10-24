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

package options

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetters(t *testing.T) {
	f := Flag{
		Long:  "test",
		Short: "t",
		Value: "hello world",
		Usage: "for testing",
	}
	assert.Equal(t, f.GetLong(), "test")
	assert.Equal(t, f.GetShort(), "t")
	assert.Equal(t, f.GetUsage(), "for testing")
}

func TestAddAll(t *testing.T) {
	cmd := &cobra.Command{}
	d, _ := time.ParseDuration("0h0m0s")
	f := flags{
		Flag{
			Long:  "int",
			Short: "i",
			Value: 1,
			Usage: "for testing",
		},
		Flag{
			Long:  "bool",
			Short: "b",
			Value: true,
			Usage: "for testing",
		},
		Flag{
			Long:  "string",
			Short: "s",
			Value: "hello world",
			Usage: "for testing",
		},
		Flag{
			Long:  "duration",
			Short: "d",
			Value: d,
			Usage: "for testing",
		},
		Flag{
			Long:  "slice",
			Short: "p",
			Value: []string{"foo"},
			Usage: "for testing",
		},
	}
	assert.Nil(t, f.AddAll(cmd))
	cobraValOne, _ := cmd.Flags().GetInt("int")
	cobraValTwo, _ := cmd.Flags().GetBool("bool")
	cobraValThree, _ := cmd.Flags().GetString("string")
	cobraValFour, _ := cmd.Flags().GetDuration("duration")
	cobraValFive, _ := cmd.Flags().GetStringSlice("slice")
	assert.Equal(t, viper.GetString("string"), "hello world")
	assert.Equal(t, viper.GetInt("int"), 1)
	assert.True(t, viper.GetBool("bool"))
	assert.Equal(t, viper.GetDuration("duration"), d)
	assert.Equal(t, cobraValOne, 1)
	assert.True(t, cobraValTwo)
	assert.Equal(t, cobraValThree, "hello world")
	assert.Equal(t, cobraValFour, d)
	assert.Equal(t, cobraValFive, []string{"foo"})
}

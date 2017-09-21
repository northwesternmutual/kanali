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

package config

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flag is a simplified representation of a configuration item.
type Flag struct {
	Long  string
	Short string
	Value interface{}
	Usage string
}

type flags []Flag

// Flags is the aggregate set of flags that Kanali has available to configure.
var Flags = &flags{}

func (f *flags) Add(a ...Flag) {
	for _, curr := range a {
		*f = append(*f, curr)
	}
}

// GetLong returns the name of the flag
func (f Flag) GetLong() string {
	return f.Long
}

// GetShort returns the Short name of the flag
func (f Flag) GetShort() string {
	return f.Short
}

// GetUsage returns the flag's description.
func (f Flag) GetUsage() string {
	return f.Usage
}

func (f flags) AddAll(cmd *cobra.Command) error {
	for _, currFlag := range f {
		switch v := currFlag.Value.(type) {
		case int:
			cmd.Flags().IntP(currFlag.Long, currFlag.Short, v, currFlag.Usage)
		case bool:
			cmd.Flags().BoolP(currFlag.Long, currFlag.Short, v, currFlag.Usage)
		case string:
			cmd.Flags().StringP(currFlag.Long, currFlag.Short, v, currFlag.Usage)
		case time.Duration:
			cmd.Flags().DurationP(currFlag.Long, currFlag.Short, v, currFlag.Usage)
		case []string:
			cmd.Flags().StringSliceP(currFlag.Long, currFlag.Short, v, currFlag.Usage)
		default:
			return nil
		}
		if err := viper.BindPFlag(currFlag.Long, cmd.Flags().Lookup(currFlag.Long)); err != nil {
			return err
		}
		viper.SetDefault(currFlag.Long, currFlag.Value)
	}

	return nil

}

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

package flags

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/northwesternmutual/kanali/pkg/log"
)

// Flag is a simplified representation of a configuration item.
type Flag struct {
	Long  string
	Short string
	Value interface{}
	Usage string
}

type flagSet []Flag

func NewFlagSet() *flagSet {
	return &flagSet{}
}

func (f *flagSet) Add(a ...Flag) {
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

func (f flagSet) AddAll(cmd *cobra.Command) error {
	for _, curr := range f {
		viper.SetDefault(curr.Long, curr.Value)
		switch v := curr.Value.(type) {
		case int:
			cmd.Flags().IntP(curr.Long, curr.Short, v, curr.Usage)
		case bool:
			cmd.Flags().BoolP(curr.Long, curr.Short, v, curr.Usage)
		case string:
			cmd.Flags().StringP(curr.Long, curr.Short, v, curr.Usage)
		case time.Duration:
			cmd.Flags().DurationP(curr.Long, curr.Short, v, curr.Usage)
		case []string:
			cmd.Flags().StringSliceP(curr.Long, curr.Short, v, curr.Usage)
		case log.Level:
			cmd.Flags().VarP(&v, curr.Long, curr.Short, curr.Usage)
		default:
			continue
		}
		if err := viper.BindPFlag(curr.Long, cmd.Flags().Lookup(curr.Long)); err != nil {
			return err
		}
	}
	return nil
}

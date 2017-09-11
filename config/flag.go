package config

import (
	"errors"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flag is a simplified representation of a configuration item.
type Flag struct {
	long  string
	short string
	value interface{}
	usage string
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
	return f.long
}

// GetShort returns the short name of the flag
func (f Flag) GetShort() string {
	return f.short
}

// GetUsage returns the flag's description.
func (f Flag) GetUsage() string {
	return f.usage
}

func (f flags) AddAll(cmd *cobra.Command) error {

	for _, currFlag := range f {
		switch v := currFlag.value.(type) {
		case int:
			cmd.Flags().IntP(currFlag.long, currFlag.short, v, currFlag.usage)
		case bool:
			cmd.Flags().BoolP(currFlag.long, currFlag.short, v, currFlag.usage)
		case string:
			cmd.Flags().StringP(currFlag.long, currFlag.short, v, currFlag.usage)
		case time.Duration:
			cmd.Flags().DurationP(currFlag.long, currFlag.short, v, currFlag.usage)
    case []string:
			cmd.Flags().StringSliceP(currFlag.long, currFlag.short, v, currFlag.usage)
		default:
			return errors.New("unsupported flag type")
		}
		if err := viper.BindPFlag(currFlag.long, cmd.Flags().Lookup(currFlag.long)); err != nil {
			return err
		}
		viper.SetDefault(currFlag.long, currFlag.value)
	}

	return nil

}

package config

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetters(t *testing.T) {
	f := flag{
		long:  "test",
		short: "t",
		value: "hello world",
		usage: "for testing",
	}
	assert.Equal(t, f.GetLong(), "test")
	assert.Equal(t, f.GetShort(), "t")
	assert.Equal(t, f.GetUsage(), "for testing")
}

func TestAddAll(t *testing.T) {
	cmd := &cobra.Command{}
	d, _ := time.ParseDuration("0h0m0s")
	f := flags{
		flag{
			long:  "int",
			short: "i",
			value: 1,
			usage: "for testing",
		},
		flag{
			long:  "bool",
			short: "b",
			value: true,
			usage: "for testing",
		},
		flag{
			long:  "string",
			short: "s",
			value: "hello world",
			usage: "for testing",
		},
		flag{
			long:  "duration",
			short: "d",
			value: d,
			usage: "for testing",
		},
	}
	assert.Nil(t, f.AddAll(cmd))
	cobraValOne, _ := cmd.Flags().GetInt("int")
	cobraValTwo, _ := cmd.Flags().GetBool("bool")
	cobraValThree, _ := cmd.Flags().GetString("string")
	cobraValFour, _ := cmd.Flags().GetDuration("duration")
	assert.Equal(t, viper.GetString("string"), "hello world")
	assert.Equal(t, viper.GetInt("int"), 1)
	assert.True(t, viper.GetBool("bool"))
	assert.Equal(t, viper.GetDuration("duration"), d)
	assert.Equal(t, cobraValOne, 1)
	assert.True(t, cobraValTwo)
	assert.Equal(t, cobraValThree, "hello world")
	assert.Equal(t, cobraValFour, d)
	f = flags{
		flag{
			long:  "wrong",
			short: "w",
			value: make(chan int),
			usage: "for testing",
		},
	}
	assert.Equal(t, f.AddAll(cmd).Error(), "unsupported flag type")
}

package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/northwesternmutual/kanali/cmd/kanali/app"
	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version string
var commit string

var rootCmd = &cobra.Command{
	Use:   "kanali",
	Short: "kubernetes native api gateway",
	Long:  "kubernetes native api gateway",
}

var startCmd = &cobra.Command{
	Use:   `start`,
	Short: `start kanali`,
	Long:  `start kanali`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.Run(context.Background()); err != nil {
			panic(err)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: `version`,
	Long:  `kanali version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("%s (%s)", version, commit))
	},
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.kanali")
	viper.AddConfigPath("/etc/kanali/")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("kanali")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.ReadInConfig()

	if err := options.KanaliOptions.AddAll(startCmd); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)

	logging.Init(nil)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

package cmd

import (
	//"fmt"
	//"os"

	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

func init() {
	// generateCmd.PersistentFlags().StringP("name", "n", "", "name of the api key")
  //
	// if err := viper.BindPFlag("name", generateCmd.PersistentFlags().Lookup("name")); err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }
  //
	// viper.SetDefault("namespace", "default")

	RootCmd.AddCommand(pluginCmd)
}

var pluginCmd = &cobra.Command{
	Use:   `plugin`,
	Short: `helper commands for Kanali plugins`,
	Long:  `helper commands for Kanali plugins`,
}

package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string
var commit string

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   `version`,
		Short: `Version information`,
		Long:  `Version information`,
		Run:   versionCmdRun,
	}
}

func versionCmdRun(cmd *cobra.Command, args []string) {
	fmt.Println(fmt.Sprintf("%s (%s)", version, commit))
}

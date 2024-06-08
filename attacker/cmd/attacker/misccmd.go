package main

import (
	"fmt"
	"github.com/tsinghua-cel/attacker-service/versions"

	"github.com/spf13/cobra"
)

var versionDetail bool

func init() {
	RootCmd.AddCommand(versionCmd)

	versionDetail = *versionCmd.Flags().BoolP("detail", "d", true, "Print detail version info")
}

// versionCmd represents the base command when called without any subcommands
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if versionDetail {
			fmt.Println(versions.DetailVersion())
		} else {
			fmt.Println(versions.Version())
		}
	},
}

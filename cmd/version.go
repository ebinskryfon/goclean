package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the goclean CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("goclean version %s\n", AppVersion)
	},
}

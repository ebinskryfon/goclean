package cmd

import (
	"fmt"
	"os"

	"github.com/ebinskryfon/goclean/pkg/ui"
	"github.com/spf13/cobra"
)

const AppVersion = "1.0.0"

var rootCmd = &cobra.Command{
	Use:     "goclean",
	Version: AppVersion,
	Short:   "goclean — Modern Go Clean Architecture generator & CLI tool",
	Long: fmt.Sprintf(`%s%s  ____      ____ _                  
 / ___|___ / ___| | ___  __ _ _ __  
| |  _/ _ \ |   | |/ _ \/ _`+"`"+` | '_ \ 
| |_| | (_) | |___| |  __/ (_| | | | |
 \____|\___/ \____|_|\___|\__,_|_| |_|  v%s%s

Go Clean Architecture Scaffolding & Code Generation Tool.
Quickly bootstrap production-ready Go projects and scaffold Clean Architecture modules.`, ui.Bold, ui.Cyan, AppVersion, ui.Reset),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(makeModuleCmd)
	rootCmd.AddCommand(makeEntityCmd)
	rootCmd.AddCommand(makeRepoCmd)
	rootCmd.AddCommand(makeUsecaseCmd)
	rootCmd.AddCommand(makeHandlerCmd)
	rootCmd.AddCommand(makeMigrationCmd)
}

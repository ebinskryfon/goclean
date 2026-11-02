package cmd

import (
	"os"

	"github.com/ebinskryfon/goclean/pkg/generator"
	"github.com/ebinskryfon/goclean/pkg/ui"
	"github.com/spf13/cobra"
)

var makeModuleCmd = &cobra.Command{
	Use:     "make:module <Name>",
	Aliases: []string{"make:resource"},
	Short:   "Generate an entire Clean Architecture module (Entity, Repo, Usecase, Handler, Migration)",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Error("Failed to determine current directory: %v", err)
			return
		}

		name := args[0]
		if err := generator.GenerateModule(cwd, name); err != nil {
			ui.Error("%v", err)
		}
	},
}

var makeEntityCmd = &cobra.Command{
	Use:   "make:entity <Name>",
	Short: "Generate a domain entity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSingleLayerGenerator("entity", args[0])
	},
}

var makeTypesCmd = &cobra.Command{
	Use:     "make:types <Name>",
	Aliases: []string{"make:type", "make:dto"},
	Short:   "Generate request, input, and response DTO types",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSingleLayerGenerator("types", args[0])
	},
}

var makeRepoCmd = &cobra.Command{
	Use:     "make:repo <Name>",
	Aliases: []string{"make:repository"},
	Short:   "Generate repository interface and GORM persistence implementation",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSingleLayerGenerator("repo", args[0])
	},
}

var makeUsecaseCmd = &cobra.Command{
	Use:   "make:usecase <Name>",
	Short: "Generate usecase interface and business logic implementation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSingleLayerGenerator("usecase", args[0])
	},
}

var makeHandlerCmd = &cobra.Command{
	Use:   "make:handler <Name>",
	Short: "Generate an HTTP Gin handler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSingleLayerGenerator("handler", args[0])
	},
}

var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration <Name>",
	Short: "Generate up and down database migration files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSingleLayerGenerator("migration", args[0])
	},
}

func runSingleLayerGenerator(layerType, name string) {
	cwd, err := os.Getwd()
	if err != nil {
		ui.Error("Failed to determine current directory: %v", err)
		return
	}

	if err := generator.GenerateSingleLayer(cwd, layerType, name); err != nil {
		ui.Error("%v", err)
	}
}

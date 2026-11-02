package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ebinskryfon/goclean/pkg/generator"
	"github.com/ebinskryfon/goclean/pkg/strutil"
	"github.com/ebinskryfon/goclean/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	flagModule string
	flagDB     string
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a fresh Go Clean Architecture project",
	Long:  "Bootstraps a complete, production-ready Go Clean Architecture backend with Gin, GORM, PostgreSQL, Viper, Zerolog, migrations, and Docker.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ui.Banner(AppVersion)

		reader := bufio.NewReader(os.Stdin)
		projectName := ""
		if len(args) > 0 {
			projectName = strings.TrimSpace(args[0])
		} else {
			fmt.Printf("  %s? Enter project name:%s ", ui.Bold, ui.Reset)
			input, _ := reader.ReadString('\n')
			projectName = strings.TrimSpace(input)
		}

		if projectName == "" {
			ui.Error("Project name cannot be empty")
			return
		}

		cleanProjectName := strutil.ToKebabCase(projectName)
		targetDir := filepath.Join(".", cleanProjectName)

		moduleName := flagModule
		if moduleName == "" {
			defaultModule := cleanProjectName
			fmt.Printf("  %s? Go module path%s [%s]: ", ui.Bold, ui.Reset, defaultModule)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input != "" {
				moduleName = input
			} else {
				moduleName = defaultModule
			}
		}

		dbName := flagDB
		if dbName == "" {
			dbName = strutil.ToSnakeCase(cleanProjectName) + "_db"
		}

		ui.Header(fmt.Sprintf("Creating project '%s'...", cleanProjectName))
		ui.Info("Module:   %s", moduleName)
		ui.Info("Database: %s", dbName)
		ui.Info("Target:   %s", targetDir)
		fmt.Println()

		opts := generator.ProjectOptions{
			ProjectName: cleanProjectName,
			ModuleName:  moduleName,
			DBName:      dbName,
			TargetDir:   targetDir,
		}

		if err := generator.GenerateProject(opts); err != nil {
			ui.Error("Failed to create project: %v", err)
			return
		}

		fmt.Println()
		ui.Success("Project '%s' created successfully!", cleanProjectName)
		fmt.Println()
		fmt.Printf("  %sNext steps:%s\n", ui.Bold, ui.Reset)
		fmt.Printf("    1. %scd %s%s\n", ui.Cyan, cleanProjectName, ui.Reset)
		fmt.Printf("    2. %sdocker compose up -d%s        # Start local PostgreSQL\n", ui.Cyan, ui.Reset)
		fmt.Printf("    3. %smake migrate%s                 # Run initial database migrations\n", ui.Cyan, ui.Reset)
		fmt.Printf("    4. %smake dev%s                     # Start server with live hot-reload\n", ui.Cyan, ui.Reset)
		fmt.Println()
		fmt.Printf("  %sScaffold new features:%s\n", ui.Bold, ui.Reset)
		fmt.Printf("    %sgoclean make:module <Name>%s     # Generates entity, repo, usecase, handler\n", ui.Cyan, ui.Reset)
		fmt.Println()
	},
}

func init() {
	newCmd.Flags().StringVarP(&flagModule, "module", "m", "", "Go module path (e.g. github.com/username/project)")
	newCmd.Flags().StringVarP(&flagDB, "db", "d", "", "PostgreSQL database name (default: <project>_db)")
}

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

type TemplateData struct {
	ProjectName string
	UseAuth     bool
	DBType      string
}

var framework string
var useDocker bool

// auth
var useAuth bool
var dbType string

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		//make folder
		err := os.Mkdir(projectName, os.ModePerm)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		data := TemplateData{
			ProjectName: projectName,
			UseAuth:     useAuth,
			DBType:      dbType,
		}

		//go mod init
		initGoMod(projectName)

		//generate main.go
		mainTemplate := "templates/main.go.tmpl"

		switch framework {
		case "fiber":
			mainTemplate = "templates/fiber/main.go.tmpl"
		case "gin":
			mainTemplate = "templates/gin/main.go.tmpl"
		}

		generateFile(
			mainTemplate,
			filepath.Join(projectName, "main.go"),
			data,
		)

		// generate README.md
		generateFile(
			"templates/readme.md.tmpl",
			filepath.Join(projectName, "README.md"),
			data,
		)
		//generate .env
		generateFile(
			"templates/env.tmpl",
			filepath.Join(projectName, ".env"),
			data,
		)

		//create structure
		createStructure(projectName)

		//docker
		if useDocker {
			generateFile(
				"templates/dockerfile.tmpl",
				filepath.Join(projectName, "Dockerfile"),
				data,
			)
		}

		//db
		if dbType == "postgres" {
			generateFile(
				"templates/db/postgres.go.tmpl",
				filepath.Join(projectName, "config/database.go"),
				data,
			)
		}

		//auth fiber
		if useAuth && framework == "fiber" {
			generateFile(
				"templates/auth/fiber/handler.go.tmpl",
				filepath.Join(projectName, "internal/handler/auth.go"),
				data,
			)
		}

		// go mod tidy
		runGoModTidy(projectName)

		fmt.Println("Project created successfully!")
	},
}

func generateFile(templatePath, outputPath string, data TemplateData) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("Template error:", err)
		return
	}

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("File error:", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Println("Execute error:", err)
	}
}

func initGoMod(projectName string) {
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = projectName

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("go mod error:", err)
		fmt.Println(string(output))
		return
	}

	fmt.Println("go.mod initialized")
}

func createStructure(projectName string) {
	dirs := []string{
		"cmd",
		"internal/handler",
		"internal/service",
		"internal/repository",
		"config",
		"pkg",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating:", path)
		}
	}
}

func runGoModTidy(projectName string) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectName

	err := cmd.Run()
	if err != nil {
		fmt.Println("go mod tidy error:", err)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&framework, "framework", "f", "net-http", "Choose framework (fiber/gin/net-http)")
	newCmd.Flags().BoolVar(&useDocker, "docker", false, "Include Docker setup")

	//auth
	newCmd.Flags().BoolVar(&useAuth, "auth", false, "Include JWT auth starter")
	newCmd.Flags().StringVar(&dbType, "db", "none", "Database (postgres/mysql/none)")
}

package cmd

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type Field struct {
	Name string
	JSON string
}

type Entity struct {
	Title  string
	Lower  string
	Plural string
	Fields []Field
}

type PkgPaths struct {
	Handler    string
	Service    string
	Repository string
	Model      string
	Route      string
}

type TemplateData struct {
	ProjectName    string
	UseAuth        bool
	DBType         string
	UseDocker      bool
	Structure      string
	WithMakefile   bool
	WithGitignore  bool
	WithMigrations bool
	HasRoutes      bool
	Entities       []Entity
	Entity         Entity
	Pkg            PkgPaths
}

var frameworks = []string{"fiber", "gin", "echo", "chi", "mux", "net-http"}
var databases = []string{"postgres", "mysql", "sqlite", "mongo", "redis", "none"}
var structures = []string{"standard", "clean", "hexagonal"}

//go:embed templates/*
var templatesFS embed.FS

var framework string
var useDocker bool

// auth
var useAuth bool
var dbType string

// crud
var crudEntities []string

// structure extras
var structure string
var withMakefile bool
var withGitignore bool
var withMigrations bool

// interactive
var interactiveMode bool

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go project",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		}

		var entities []Entity
		if interactiveMode || projectName == "" {
			var err error
			projectName, entities, err = interactiveSetup(projectName)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			if projectName == "" {
				fmt.Println("Error: project name is required")
				return
			}
		} else {
			entities = parseEntities(crudEntities)
		}

		if !validChoice(frameworks, framework) {
			framework = "net-http"
		}
		if !validChoice(databases, dbType) {
			dbType = "none"
		}
		if !validChoice(structures, structure) {
			structure = "standard"
		}

		err := os.Mkdir(projectName, os.ModePerm)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		data := TemplateData{
			ProjectName:    projectName,
			UseAuth:        useAuth,
			DBType:         dbType,
			UseDocker:      useDocker,
			Structure:      structure,
			WithMakefile:   withMakefile,
			WithGitignore:  withGitignore,
			WithMigrations: withMigrations,
			Entities:       entities,
			HasRoutes:      useAuth || len(entities) > 0,
			Pkg:            pathsForStructure(structure),
		}

		//go mod init
		initGoMod(projectName)

		//generate main.go
		generateFile(
			"templates/"+frameworkDir(framework)+"/main.go.tmpl",
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
		createStructure(projectName, structure)

		//docker
		if useDocker {
			generateFile(
				"templates/dockerfile.tmpl",
				filepath.Join(projectName, "Dockerfile"),
				data,
			)
		}

		//db
		if dbType != "none" {
			generateFile(
				"templates/db/"+dbType+".go.tmpl",
				filepath.Join(projectName, "config/database.go"),
				data,
			)
		}

		//auth
		if useAuth {
			generateAuth(framework, projectName, data)
		}

		//crud
		if len(entities) > 0 {
			generateCRUD(framework, projectName, data)
		}

		//route
		if data.HasRoutes {
			generateRoute(framework, projectName, data)
		}

		//http-style shared helpers
		if (useAuth || len(entities) > 0) && handlerTemplateDir(framework) == "http" {
			generateFile(
				"templates/crud/helpers.go.tmpl",
				filepath.Join(projectName, data.Pkg.Handler, "helpers.go"),
				data,
			)
		}

		//extras
		if withMakefile {
			generateFile(
				"templates/makefile.tmpl",
				filepath.Join(projectName, "Makefile"),
				data,
			)
		}
		if withGitignore {
			generateFile(
				"templates/gitignore.tmpl",
				filepath.Join(projectName, ".gitignore"),
				data,
			)
		}
		if withMigrations {
			generateFile(
				"templates/migration.sql.tmpl",
				filepath.Join(projectName, "migrations/0001_init.sql"),
				data,
			)
		}

		// go mod tidy
		runGoModTidy(projectName)

		// gofmt
		runGoFmt(projectName)

		fmt.Println("Project created successfully!")
	},
}

func frameworkDir(fw string) string {
	if fw == "net-http" {
		return "http"
	}
	return fw
}

func handlerTemplateDir(fw string) string {
	switch fw {
	case "net-http", "http", "chi", "mux":
		return "http"
	}
	return fw
}

func validChoice(choices []string, value string) bool {
	for _, c := range choices {
		if c == value {
			return true
		}
	}
	return false
}

func pathsForStructure(s string) PkgPaths {
	switch s {
	case "clean":
		return PkgPaths{
			Handler:    "internal/delivery/http",
			Service:    "internal/usecases",
			Repository: "internal/repositories",
			Model:      "internal/entities",
			Route:      "internal/route",
		}
	case "hexagonal":
		return PkgPaths{
			Handler:    "internal/adapters/in/http",
			Service:    "internal/core/service",
			Repository: "internal/adapters/out",
			Model:      "internal/core/domain",
			Route:      "internal/adapters/in/http/route",
		}
	default:
		return PkgPaths{
			Handler:    "internal/handler",
			Service:    "internal/service",
			Repository: "internal/repository",
			Model:      "internal/model",
			Route:      "internal/route",
		}
	}
}

func parseEntities(inputs []string) []Entity {
	var entities []Entity
	for _, in := range inputs {
		entity, ok := parseEntity(in)
		if ok {
			entities = append(entities, entity)
		}
	}
	return entities
}

func parseEntity(input string) (Entity, bool) {
	parts := strings.SplitN(input, ":", 2)
	name := strings.TrimSpace(parts[0])
	if name == "" {
		return Entity{}, false
	}

	lower := strings.ToLower(name)
	title := capitalize(lower)

	var fields []Field
	if len(parts) > 1 {
		for _, f := range strings.Split(parts[1], ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				fields = append(fields, Field{Name: capitalize(f), JSON: strings.ToLower(f)})
			}
		}
	}
	if len(fields) == 0 {
		fields = []Field{
			{Name: "Name", JSON: "name"},
			{Name: "Email", JSON: "email"},
		}
	}

	return Entity{
		Title:  title,
		Lower:  lower,
		Plural: pluralize(lower),
		Fields: fields,
	}, true
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func pluralize(s string) string {
	if strings.HasSuffix(s, "y") && len(s) > 1 {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") ||
		strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}
	return s + "s"
}

func generateFile(templatePath, outputPath string, data TemplateData) {
	tmpl, err := template.ParseFS(templatesFS, templatePath)
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

func generateAuth(fw, projectName string, data TemplateData) {
	generateFile(
		"templates/auth/"+handlerTemplateDir(fw)+"/handler.go.tmpl",
		filepath.Join(projectName, data.Pkg.Handler, "auth.go"),
		data,
	)
}

func generateCRUD(fw, projectName string, data TemplateData) {
	dir := handlerTemplateDir(fw)

	for _, e := range data.Entities {
		ed := data
		ed.Entity = e

		generateFile(
			"templates/crud/model.go.tmpl",
			filepath.Join(projectName, data.Pkg.Model, e.Lower+".go"),
			ed,
		)
		generateFile(
			"templates/crud/repository.go.tmpl",
			filepath.Join(projectName, data.Pkg.Repository, e.Lower+"_repo.go"),
			ed,
		)
		generateFile(
			"templates/crud/service.go.tmpl",
			filepath.Join(projectName, data.Pkg.Service, e.Lower+"_service.go"),
			ed,
		)
		generateFile(
			"templates/crud/"+dir+"/handler.go.tmpl",
			filepath.Join(projectName, data.Pkg.Handler, e.Lower+".go"),
			ed,
		)
	}

	// ErrNotFound shared by all services
	generateFile(
		"templates/crud/errors.go.tmpl",
		filepath.Join(projectName, data.Pkg.Service, "errors.go"),
		data,
	)

	// writeJSON / idFromPath shared by net/http-style handlers
	if dir == "http" {
		generateFile(
			"templates/crud/helpers.go.tmpl",
			filepath.Join(projectName, data.Pkg.Handler, "helpers.go"),
			data,
		)
	}
}

func generateRoute(fw, projectName string, data TemplateData) {
	generateFile(
		"templates/route/"+frameworkDir(fw)+"/route.go.tmpl",
		filepath.Join(projectName, data.Pkg.Route, "route.go"),
		data,
	)
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

func createStructure(projectName, structure string) {
	var dirs []string

	switch structure {
	case "clean":
		dirs = []string{
			"cmd",
			"internal/entities",
			"internal/usecases",
			"internal/repositories",
			"internal/delivery/http",
			"internal/route",
			"config",
			"pkg",
		}
	case "hexagonal":
		dirs = []string{
			"cmd",
			"internal/core/domain",
			"internal/core/service",
			"internal/adapters/out",
			"internal/adapters/in/http/route",
			"config",
			"pkg",
		}
	default:
		dirs = []string{
			"cmd",
			"internal/model",
			"internal/handler",
			"internal/service",
			"internal/repository",
			"internal/route",
			"config",
			"pkg",
		}
	}

	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating:", path)
		}
	}

	if withMigrations {
		err := os.MkdirAll(filepath.Join(projectName, "migrations"), os.ModePerm)
		if err != nil {
			fmt.Println("Error creating migrations dir")
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

func runGoFmt(projectName string) {
	cmd := exec.Command("gofmt", "-w", ".")
	cmd.Dir = projectName

	err := cmd.Run()
	if err != nil {
		fmt.Println("gofmt error:", err)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&framework, "framework", "f", "net-http", "Choose framework (fiber/gin/echo/chi/mux/net-http)")
	newCmd.Flags().BoolVar(&useDocker, "docker", false, "Include Docker setup")

	//auth
	newCmd.Flags().BoolVar(&useAuth, "auth", false, "Include JWT auth starter")
	newCmd.Flags().StringVar(&dbType, "db", "none", "Database (postgres/mysql/sqlite/mongo/redis/none)")

	//crud
	newCmd.Flags().StringArrayVar(&crudEntities, "crud", nil, "Generate CRUD for entity, format: --crud=user:name,email (repeatable)")

	//structure
	newCmd.Flags().StringVar(&structure, "structure", "standard", "Project layout (standard/clean/hexagonal)")
	newCmd.Flags().BoolVar(&withMakefile, "makefile", false, "Include Makefile")
	newCmd.Flags().BoolVar(&withGitignore, "gitignore", false, "Include .gitignore")
	newCmd.Flags().BoolVar(&withMigrations, "migrations", false, "Include migrations folder")

	//interactive
	newCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "Run in interactive mode")
}

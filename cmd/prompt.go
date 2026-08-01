package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var stdinReader = bufio.NewReader(os.Stdin)

func interactiveSetup(defaultName string) (string, []Entity, error) {
	name := defaultName
	if name == "" {
		name = prompt("Project name", "")
		if name == "" {
			return "", nil, errors.New("project name is required")
		}
	}

	framework = prompt("Framework ("+strings.Join(frameworks, "/")+")", "net-http")
	if !validChoice(frameworks, framework) {
		return "", nil, errors.New("unknown framework: " + framework)
	}

	dbType = prompt("Database ("+strings.Join(databases, "/")+")", "none")
	if !validChoice(databases, dbType) {
		return "", nil, errors.New("unknown database: " + dbType)
	}

	useAuth = promptYesNo("Include JWT auth?")
	useDocker = promptYesNo("Include Docker support?")

	structure = prompt("Structure ("+strings.Join(structures, "/")+")", "standard")
	if !validChoice(structures, structure) {
		return "", nil, errors.New("unknown structure: " + structure)
	}

	withMakefile = promptYesNo("Include Makefile?")
	withGitignore = promptYesNo("Include .gitignore?")
	withMigrations = promptYesNo("Include migrations folder?")

	var entities []Entity
	for {
		entityName := prompt("CRUD entity name (blank to skip)", "")
		if entityName == "" {
			break
		}
		fields := prompt("Fields, comma-separated (blank for default: name,email)", "")
		entity, ok := parseEntity(entityName + ":" + fields)
		if !ok {
			continue
		}
		entities = append(entities, entity)
	}

	return name, entities, nil
}

func prompt(question, def string) string {
	if def != "" {
		fmt.Printf("%s [%s]: ", question, def)
	} else {
		fmt.Printf("%s: ", question)
	}

	line, err := stdinReader.ReadString('\n')
	if err != nil && line == "" {
		return def
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}

func promptYesNo(question string) bool {
	for {
		answer := strings.ToLower(prompt(question+" (y/n)", "n"))
		switch answer {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		}
	}
}

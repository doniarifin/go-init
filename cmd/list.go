package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available options",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Frameworks: " + strings.Join(frameworks, ", "))
		fmt.Println("Databases:  " + strings.Join(databases, ", "))
		fmt.Println("Structures: " + strings.Join(structures, ", "))
		fmt.Println()
		fmt.Println("Flags for 'go-init new':")
		fmt.Println("  --framework, -f   web framework")
		fmt.Println("  --db              database driver")
		fmt.Println("  --auth            JWT auth starter")
		fmt.Println("  --crud=entity     CRUD boilerplate, format: --crud=user:name,email (repeatable)")
		fmt.Println("  --docker          Dockerfile")
		fmt.Println("  --structure       project layout (standard/clean/hexagonal)")
		fmt.Println("  --makefile        Makefile")
		fmt.Println("  --gitignore       .gitignore")
		fmt.Println("  --migrations      migrations folder")
		fmt.Println("  --interactive, -i guided setup")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

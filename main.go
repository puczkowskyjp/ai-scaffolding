package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/puczkowskyjp/ai-scaffolding/internal/generation"
	"github.com/puczkowskyjp/ai-scaffolding/internal/initcmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ai-scaffold <command> [options] [path]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  init    Initialize AI configuration")
		return
	}

	root := "."

	switch os.Args[1] {
	case "init":
		initFlags := flag.NewFlagSet("init", flag.ContinueOnError)
		initFlags.SetOutput(os.Stdout)

		nonInteractive := initFlags.Bool("non-interactive", false, "Generate scaffolding without interactive prompts")
		if err := initFlags.Parse(os.Args[2:]); err != nil {
			os.Exit(1)
		}

		args := initFlags.Args()
		if len(args) > 1 {
			fmt.Println("Usage: ai-scaffold init [--non-interactive] [path]")
			os.Exit(1)
		}

		if len(args) == 1 {
			root = args[0]
		}

		templateRoot, err := generation.GetTemplateRoot()
		if err != nil {
			fmt.Println("Error locating templates:", err)
			os.Exit(1)
		}

		err = initcmd.Run(initcmd.Options{
			Root:           root,
			TemplateRoot:   templateRoot,
			NonInteractive: *nonInteractive,
			Input:          os.Stdin,
			Output:         os.Stdout,
		})
		if err != nil {
			if errors.Is(err, initcmd.ErrCancelled) {
				fmt.Println()
				fmt.Println("Generation cancelled.")
				return
			}

			fmt.Println("Error generating configuration:", err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Println("AI configuration generated.")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

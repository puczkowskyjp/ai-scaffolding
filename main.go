package main

import (
	"fmt"
	"os"

	"github.com/puczkowskyjp/ai-scaffolding/internal/detection"
	"github.com/puczkowskyjp/ai-scaffolding/internal/generation"
	"github.com/puczkowskyjp/ai-scaffolding/internal/planning"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ai-scaffold <command>")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  init    Initialize AI configuration")
		return
	}

	root := "."
	if len(os.Args) >= 3 {
		root = os.Args[2]
	}

	switch os.Args[1] {
	case "init":
		fmt.Println("Initializing AI configuration...")

		profile := detection.Detect(root)

		fmt.Println()
		fmt.Println("Architecture:")
		fmt.Println("  1. Monolith")
		fmt.Println("  2. Modular monolith")
		fmt.Println("  3. Microservices")
		fmt.Println("  4. Serverless")
		fmt.Println("  5. Not sure")

		fmt.Print("Select architecture: ")

		var architectureChoice int
		fmt.Scanln(&architectureChoice)

		var architecture detection.Architecture

		switch architectureChoice {
		case 1:
			architecture = detection.ArchitectureMonolith
		case 2:
			architecture = detection.ArchitectureModularMonolith
		case 3:
			architecture = detection.ArchitectureMicroservices
		case 4:
			architecture = detection.ArchitectureServerless
		default:
			architecture = detection.ArchitectureUnknown
		}

		profile.Architecture = architecture

		fmt.Print("Generate testing agents? (y/n): ")

		var answer string
		fmt.Scanln(&answer)

		generateTesting := answer == "y" || answer == "Y"

		generateAdversary := false

		if generateTesting {
			fmt.Print("Generate adversary agent? (y/n): ")
			fmt.Scanln(&answer)

			generateAdversary = answer == "y" || answer == "Y"
			fmt.Println(generateAdversary)
		}

		plan := planning.BuildPlan(
			profile,
			generateTesting,
			generateAdversary,
		)

		fmt.Println()
		fmt.Println("Recommended agents:")

		templateRoot, err := generation.GetTemplateRoot()
		if err != nil {
			fmt.Println("Error locating templates:", err)
			return
		}

		if err := generation.Generate(
			root,
			templateRoot,
			plan,
			profile,
		); err != nil {
			fmt.Println("Error generating configuration:", err)
			return
		}

		fmt.Println()
		fmt.Println("AI configuration generated.")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}

}

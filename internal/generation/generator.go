package generation

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/puczkowskyjp/ai-scaffolding/internal/detection"
	"github.com/puczkowskyjp/ai-scaffolding/internal/planning"
)

// GenerationContext bundles the data that is available to agent and instruction
// templates during rendering.
type GenerationContext struct {
	Profile detection.ProjectProfile
	Plan    planning.AgentPlan
}

// PlannedFiles returns the relative file paths that Generate will attempt to
// create for the provided plan.
func PlannedFiles(plan planning.AgentPlan) []string {
	files := []string{
		filepath.Join(".github", "copilot-instructions.md"),
	}

	for _, agent := range plan.Agents {
		files = append(files, filepath.Join(".github", "agents", getTemplateName(agent)))
	}

	return files
}

// Generate renders agent templates and the copilot-instructions file into the
// target repository at root. It uses templateRoot as the base directory for
// template files and skips any agent whose template is not found.
func Generate(
	root string,
	templateRoot string,
	plan planning.AgentPlan,
	profile detection.ProjectProfile,
) error {
	agentsDir := filepath.Join(root, ".github", "agents")

	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return err
	}

	context := GenerationContext{
		Profile: profile,
		Plan:    plan,
	}

	agentsTemplateRoot := filepath.Join(templateRoot, "agents")

	for _, agent := range plan.Agents {
		templateName := getTemplateName(agent)

		templatePath := filepath.Join(agentsTemplateRoot, templateName)
		targetPath := filepath.Join(agentsDir, templateName)

		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Skipped:", templateName, "(template not found)")
				continue
			}

			return err
		}

		file, err := os.Create(targetPath)
		if err != nil {
			return err
		}

		err = tmpl.Execute(file, context)
		file.Close()

		if err != nil {
			return err
		}

		fmt.Println("Generated:", templateName)
	}

	if err := generateInstructions(root, templateRoot, context); err != nil {
		return err
	}

	return nil
}

// GetTemplateRoot returns the absolute path of the templates directory relative
// to the current working directory.
func GetTemplateRoot() (string, error) {
	return filepath.Abs("templates")
}

func getTemplateName(agent planning.Agent) string {
	return agent.FileName
}

func generateInstructions(
	root string,
	templateRoot string,
	context GenerationContext,
) error {
	templatePath := filepath.Join(templateRoot, "copilot-instructions.md")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	targetPath := filepath.Join(root, ".github", "copilot-instructions.md")

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, context)
}

package initcmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/puczkowskyjp/ai-scaffolding/internal/detection"
	"github.com/puczkowskyjp/ai-scaffolding/internal/generation"
	"github.com/puczkowskyjp/ai-scaffolding/internal/planning"
)

// ErrCancelled indicates the user cancelled the workflow before generation.
var ErrCancelled = errors.New("generation cancelled")

// Options configures the init workflow.
type Options struct {
	Root           string
	TemplateRoot   string
	NonInteractive bool
	Input          io.Reader
	Output         io.Writer
}

type plannerOptions struct {
	Architecture      detection.Architecture
	GenerateTesting   bool
	GenerateAdversary bool
}

// Run executes the init workflow.
func Run(options Options) error {
	input := options.Input
	if input == nil {
		input = strings.NewReader("")
	}

	output := options.Output
	if output == nil {
		output = io.Discard
	}

	fmt.Fprintln(output, "AI Scaffold")
	fmt.Fprintln(output, "─────────────────────────────────────")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Analyzing project...")

	profile, err := detection.Detect(options.Root)
	if err != nil {
		return fmt.Errorf("discover project: %w", err)
	}

	defaultOptions := plannerOptions{
		Architecture: profile.Architecture,
	}

	if options.NonInteractive {
		profile.Architecture = defaultOptions.Architecture
		plan := planning.BuildPlan(profile, defaultOptions.GenerateTesting, defaultOptions.GenerateAdversary)
		return generation.Generate(options.Root, options.TemplateRoot, plan, profile)
	}

	session := terminalSession{
		reader: bufio.NewReader(input),
		writer: output,
	}

	for {
		session.printSummary(profile)

		action, err := session.promptDiscoveryAction()
		if err != nil {
			return err
		}

		if action == discoveryActionCancel {
			return ErrCancelled
		}

		selectedOptions := defaultOptions
		if action == discoveryActionCustomize {
			selectedOptions, err = session.promptPlannerOptions(defaultOptions)
			if err != nil {
				return err
			}
		}

		profile.Architecture = selectedOptions.Architecture
		plan := planning.BuildPlan(profile, selectedOptions.GenerateTesting, selectedOptions.GenerateAdversary)

		plan, err = session.selectAgents(plan)
		if err != nil {
			return err
		}

		session.printPreview(generation.PlannedFiles(plan))

		confirmation, err := session.promptConfirmation()
		if err != nil {
			return err
		}

		switch confirmation {
		case confirmationGenerate:
			return generation.Generate(options.Root, options.TemplateRoot, plan, profile)
		case confirmationBack:
			continue
		default:
			return ErrCancelled
		}
	}
}

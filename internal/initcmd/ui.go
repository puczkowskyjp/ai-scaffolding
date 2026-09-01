package initcmd

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/puczkowskyjp/ai-scaffolding/internal/detection"
	"github.com/puczkowskyjp/ai-scaffolding/internal/planning"
)

type discoveryAction string

const (
	discoveryActionContinue  discoveryAction = "continue"
	discoveryActionCustomize discoveryAction = "customize"
	discoveryActionCancel    discoveryAction = "cancel"
)

type confirmationAction string

const (
	confirmationGenerate confirmationAction = "generate"
	confirmationBack     confirmationAction = "back"
	confirmationCancel   confirmationAction = "cancel"
)

type terminalSession struct {
	reader lineReader
	writer io.Writer
}

type lineReader interface {
	ReadString(delim byte) (string, error)
}

func (session terminalSession) printSummary(profile detection.ProjectProfile) {
	fmt.Fprintln(session.writer)
	fmt.Fprintln(session.writer, "Project")
	fmt.Fprintf(session.writer, "  Name:           %s\n", profile.Name)
	fmt.Fprintf(session.writer, "  Languages:      %s\n", joinOrDefault(profile.Languages))
	fmt.Fprintf(session.writer, "  Frameworks:     %s\n", joinOrDefault(profile.Frameworks))
	fmt.Fprintf(session.writer, "  Project types:  %s\n", joinOrDefault(profile.ProjectTypes))
	fmt.Fprintf(session.writer, "  Architecture:   %s\n", profile.Architecture)
	fmt.Fprintln(session.writer)
	fmt.Fprintln(session.writer, "Detected Architecture")
	for _, characteristic := range profile.Characteristics {
		fmt.Fprintf(session.writer, "  ✓ %s\n", characteristic)
	}

	if len(profile.Characteristics) == 0 {
		fmt.Fprintln(session.writer, "  ✓ No specific characteristics detected")
	}
}

func (session terminalSession) promptDiscoveryAction() (discoveryAction, error) {
	for {
		fmt.Fprintln(session.writer)
		fmt.Fprint(session.writer, "[C]ontinue  [U]stomize  Cance[l]: ")

		answer, err := session.readTrimmedLine()
		if err != nil {
			return "", err
		}

		switch strings.ToLower(answer) {
		case "", "c", "continue":
			return discoveryActionContinue, nil
		case "u", "customize":
			return discoveryActionCustomize, nil
		case "l", "cancel":
			return discoveryActionCancel, nil
		}

		fmt.Fprintln(session.writer, "Please enter continue, customize, or cancel.")
	}
}

func (session terminalSession) promptPlannerOptions(defaults plannerOptions) (plannerOptions, error) {
	options := defaults

	fmt.Fprintln(session.writer)
	fmt.Fprintln(session.writer, "Architecture")
	fmt.Fprintln(session.writer, "  1. Monolith")
	fmt.Fprintln(session.writer, "  2. Modular monolith")
	fmt.Fprintln(session.writer, "  3. Microservices")
	fmt.Fprintln(session.writer, "  4. Serverless")
	fmt.Fprintln(session.writer, "  5. Unknown")

	for {
		fmt.Fprint(session.writer, "Select architecture [5]: ")
		answer, err := session.readTrimmedLine()
		if err != nil {
			return plannerOptions{}, err
		}

		if answer == "" || answer == "5" {
			options.Architecture = detection.ArchitectureUnknown
			break
		}

		switch answer {
		case "1":
			options.Architecture = detection.ArchitectureMonolith
		case "2":
			options.Architecture = detection.ArchitectureModularMonolith
		case "3":
			options.Architecture = detection.ArchitectureMicroservices
		case "4":
			options.Architecture = detection.ArchitectureServerless
		default:
			fmt.Fprintln(session.writer, "Please enter a number between 1 and 5.")
			continue
		}

		break
	}

	testing, err := session.promptYesNo("Generate testing agents? [y/N]: ")
	if err != nil {
		return plannerOptions{}, err
	}

	options.GenerateTesting = testing
	options.GenerateAdversary = false

	if options.GenerateTesting {
		adversary, err := session.promptYesNo("Generate adversary agents? [y/N]: ")
		if err != nil {
			return plannerOptions{}, err
		}

		options.GenerateAdversary = adversary
	}

	return options, nil
}

func (session terminalSession) selectAgents(plan planning.AgentPlan) (planning.AgentPlan, error) {
	enabled := make([]bool, len(plan.Agents))
	for index := range enabled {
		enabled[index] = true
	}

	for {
		fmt.Fprintln(session.writer)
		fmt.Fprintln(session.writer, "Agents to generate")

		for index, agent := range plan.Agents {
			marker := " "
			if enabled[index] {
				marker = "x"
			}

			fmt.Fprintf(session.writer, "  [%s] %d. %s\n", marker, index+1, title(agent.Name))
		}

		fmt.Fprint(session.writer, "Enter agent numbers to toggle (comma-separated), or press Enter to continue: ")

		answer, err := session.readTrimmedLine()
		if err != nil {
			return planning.AgentPlan{}, err
		}

		if answer == "" {
			return planning.ApplySelection(plan, enabled)
		}

		indices, err := parseToggleList(answer, len(plan.Agents))
		if err != nil {
			fmt.Fprintf(session.writer, "Invalid selection: %v\n", err)
			continue
		}

		for _, index := range indices {
			enabled[index] = !enabled[index]
		}
	}
}

func (session terminalSession) printPreview(paths []string) {
	fmt.Fprintln(session.writer)
	fmt.Fprintln(session.writer, "Files to generate")

	for _, line := range FormatPreviewTree(paths) {
		fmt.Fprintf(session.writer, "  %s\n", line)
	}
}

func (session terminalSession) promptConfirmation() (confirmationAction, error) {
	for {
		fmt.Fprintln(session.writer)
		fmt.Fprint(session.writer, "[G]enerate  [B]ack  Cance[l]: ")

		answer, err := session.readTrimmedLine()
		if err != nil {
			return "", err
		}

		switch strings.ToLower(answer) {
		case "", "g", "generate":
			return confirmationGenerate, nil
		case "b", "back":
			return confirmationBack, nil
		case "l", "cancel":
			return confirmationCancel, nil
		}

		fmt.Fprintln(session.writer, "Please enter generate, back, or cancel.")
	}
}

func (session terminalSession) promptYesNo(prompt string) (bool, error) {
	for {
		fmt.Fprint(session.writer, prompt)

		answer, err := session.readTrimmedLine()
		if err != nil {
			return false, err
		}

		switch strings.ToLower(answer) {
		case "y", "yes":
			return true, nil
		case "", "n", "no":
			return false, nil
		}

		fmt.Fprintln(session.writer, "Please enter y or n.")
	}
}

func (session terminalSession) readTrimmedLine() (string, error) {
	line, err := session.reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func parseToggleList(value string, max int) ([]int, error) {
	parts := strings.Split(value, ",")
	indices := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		index, err := strconv.Atoi(part)
		if err != nil || index < 1 || index > max {
			return nil, fmt.Errorf("expected values between 1 and %d", max)
		}

		indices = append(indices, index-1)
	}

	slices.Sort(indices)
	return slices.Compact(indices), nil
}

func joinOrDefault(values []string) string {
	if len(values) == 0 {
		return "None detected"
	}

	return strings.Join(values, ", ")
}

func title(name string) string {
	parts := strings.Split(name, "-")
	for index, part := range parts {
		if part == "" {
			continue
		}

		parts[index] = strings.ToUpper(part[:1]) + part[1:]
	}

	return strings.Join(parts, " ")
}

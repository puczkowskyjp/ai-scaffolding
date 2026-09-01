package initcmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunNonInteractiveGeneratesDefaultPlan(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "app.csproj"), "<Project />")

	templateRoot := createTemplates(t)
	output := &bytes.Buffer{}

	err := Run(Options{
		Root:           root,
		TemplateRoot:   templateRoot,
		NonInteractive: true,
		Output:         output,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	assertExists(t, filepath.Join(root, ".github", "copilot-instructions.md"))
	assertExists(t, filepath.Join(root, ".github", "agents", "architect.agent.md"))
	assertExists(t, filepath.Join(root, ".github", "agents", "dotnet.agent.md"))
}

func TestRunCancelDoesNotGenerateFiles(t *testing.T) {
	root := t.TempDir()
	templateRoot := createTemplates(t)

	err := Run(Options{
		Root:         root,
		TemplateRoot: templateRoot,
		Input:        strings.NewReader("l\n"),
		Output:       &bytes.Buffer{},
	})
	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}

	if _, statErr := os.Stat(filepath.Join(root, ".github")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no generated files, stat err: %v", statErr)
	}
}

func TestRunInteractiveSelectionRemovesAgentFromGeneratedFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "app.csproj"), "<Project />")
	writeFile(t, filepath.Join(root, "web.tsx"), "export const App = () => null;")

	templateRoot := createTemplates(t)

	err := Run(Options{
		Root:         root,
		TemplateRoot: templateRoot,
		Input:        strings.NewReader("c\n3\n\ng\n"),
		Output:       &bytes.Buffer{},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	assertExists(t, filepath.Join(root, ".github", "agents", "architect.agent.md"))
	assertExists(t, filepath.Join(root, ".github", "agents", "dotnet.agent.md"))

	if _, statErr := os.Stat(filepath.Join(root, ".github", "agents", "frontend.agent.md")); !os.IsNotExist(statErr) {
		t.Fatalf("expected frontend agent to be omitted, stat err: %v", statErr)
	}
}

func TestRunBackThenCancelLeavesRepositoryUnchanged(t *testing.T) {
	root := t.TempDir()
	templateRoot := createTemplates(t)

	err := Run(Options{
		Root:         root,
		TemplateRoot: templateRoot,
		Input:        strings.NewReader("c\n\nb\nl\n"),
		Output:       &bytes.Buffer{},
	})
	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}

	if _, statErr := os.Stat(filepath.Join(root, ".github")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no generated files, stat err: %v", statErr)
	}
}

func TestRunReportsInvalidTargetDirectory(t *testing.T) {
	err := Run(Options{
		Root:         filepath.Join(t.TempDir(), "missing"),
		TemplateRoot: createTemplates(t),
		Output:       &bytes.Buffer{},
	})
	if err == nil || !strings.Contains(err.Error(), "discover project") {
		t.Fatalf("expected discovery error, got %v", err)
	}
}

func TestFormatPreviewTree(t *testing.T) {
	lines := FormatPreviewTree([]string{
		filepath.Join(".github", "copilot-instructions.md"),
		filepath.Join(".github", "agents", "architect.agent.md"),
		filepath.Join(".github", "agents", "frontend.agent.md"),
	})

	expected := []string{
		".github/",
		"├── agents/",
		"│   ├── architect.agent.md",
		"│   └── frontend.agent.md",
		"└── copilot-instructions.md",
	}

	if strings.Join(lines, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("unexpected preview tree:\nexpected:\n%s\nactual:\n%s", strings.Join(expected, "\n"), strings.Join(lines, "\n"))
	}
}

func createTemplates(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "copilot-instructions.md"), "instructions {{len .Plan.Agents}}")
	writeFile(t, filepath.Join(root, "agents", "architect.agent.md"), "architect")
	writeFile(t, filepath.Join(root, "agents", "dotnet.agent.md"), "dotnet")
	writeFile(t, filepath.Join(root, "agents", "frontend.agent.md"), "frontend")

	return root
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

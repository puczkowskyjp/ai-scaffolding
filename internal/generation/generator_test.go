package generation

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/puczkowskyjp/ai-scaffolding/internal/planning"
)

func TestPlannedFilesMatchesGenerationPlan(t *testing.T) {
	plan := planning.AgentPlan{
		Agents: []planning.Agent{
			{FileName: "architect.agent.md"},
			{FileName: "frontend.agent.md"},
		},
	}

	expected := []string{
		filepath.Join(".github", "copilot-instructions.md"),
		filepath.Join(".github", "agents", "architect.agent.md"),
		filepath.Join(".github", "agents", "frontend.agent.md"),
	}

	if actual := PlannedFiles(plan); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unexpected planned files:\nexpected: %#v\nactual:   %#v", expected, actual)
	}
}

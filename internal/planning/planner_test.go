package planning

import "testing"

func TestApplySelectionRemovesDisabledAgents(t *testing.T) {
	plan := AgentPlan{
		Agents: []Agent{
			{Name: AgentArchitect},
			{Name: AgentDotNet},
			{Name: AgentBackendTesting},
		},
	}

	selected, err := ApplySelection(plan, []bool{true, false, true})
	if err != nil {
		t.Fatalf("ApplySelection returned error: %v", err)
	}

	if len(selected.Agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(selected.Agents))
	}

	if selected.Agents[0].Name != AgentArchitect || selected.Agents[1].Name != AgentBackendTesting {
		t.Fatalf("unexpected selection result: %#v", selected.Agents)
	}
}

func TestApplySelectionRejectsMismatchedState(t *testing.T) {
	_, err := ApplySelection(AgentPlan{Agents: []Agent{{Name: AgentArchitect}}}, []bool{})
	if err != ErrInvalidSelection {
		t.Fatalf("expected ErrInvalidSelection, got %v", err)
	}
}

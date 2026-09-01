package planning

import "github.com/puczkowskyjp/ai-scaffolding/internal/detection"

const (
	// AgentArchitect is the name of the architect agent that coordinates
	// design discussions and delegates to specialist agents.
	AgentArchitect = "architect"
	// AgentDotNet is the name of the agent responsible for .NET backend development.
	AgentDotNet = "dotnet"
	// AgentFrontend is the name of the agent responsible for React frontend development.
	AgentFrontend = "frontend"
	// AgentDatabase is the name of the agent responsible for database development and data access.
	AgentDatabase = "database"
	// AgentBackendTesting is the name of the agent responsible for backend testing.
	AgentBackendTesting = "backend-testing"
	// AgentFrontendTesting is the name of the agent responsible for frontend testing.
	AgentFrontendTesting = "frontend-testing"
	// AgentBackendAdversary is the name of the agent that challenges backend implementations.
	AgentBackendAdversary = "backend-adversary"
	// AgentFrontendAdversary is the name of the agent that challenges frontend implementations.
	AgentFrontendAdversary = "frontend-adversary"
	// AgentMicroservices is the name of the agent responsible for distributed service architecture.
	AgentMicroservices = "microservices"
)

// AgentPlan holds the ordered list of agents that should be generated for a
// project based on its detected profile and the user's preferences.
type AgentPlan struct {
	Agents []Agent
}

// BuildPlan constructs an AgentPlan based on the detected project profile and
// the user's choices about testing and adversary agents. The returned plan
// lists agents in the order they should be generated.
func BuildPlan(
	profile detection.ProjectProfile,
	generateTesting bool,
	generateAdversary bool,
) AgentPlan {
	plan := AgentPlan{}

	if profile.Architecture == detection.ArchitectureMicroservices {
		plan.Agents = append(plan.Agents, Agent{
			Name:        AgentMicroservices,
			FileName:    "microservice.agent.md",
			Description: "Handles distributed service architecture.",
		})
	}

	plan.Agents = append(plan.Agents, Agent{
		Name:        AgentArchitect,
		FileName:    "architect.agent.md",
		Description: "Coordinates architecture and specialist agents.",
	})

	if profile.IsDotNet {
		plan.Agents = append(plan.Agents, Agent{
			Name:        AgentDotNet,
			FileName:    "dotnet.agent.md",
			Description: "Handles .NET backend development.",
		})
	}

	if profile.IsReact {
		plan.Agents = append(plan.Agents, Agent{
			Name:        AgentFrontend,
			FileName:    "frontend.agent.md",
			Description: "Handles React frontend development.",
		})
	}

	if profile.IsPostgres {
		plan.Agents = append(plan.Agents, Agent{
			Name:        AgentDatabase,
			FileName:    "database.agent.md",
			Description: "Handles database development and data access.",
		})
	}

	if generateTesting {
		if profile.IsDotNet {
			plan.Agents = append(plan.Agents, Agent{
				Name:        AgentBackendTesting,
				FileName:    "backend-testing.agent.md",
				Description: "Handles backend testing.",
			})
		}

		if profile.IsReact {
			plan.Agents = append(plan.Agents, Agent{
				Name:        AgentFrontendTesting,
				FileName:    "frontend-testing.agent.md",
				Description: "Handles frontend testing.",
			})
		}
	}

	if generateAdversary {
		if profile.IsDotNet {
			plan.Agents = append(plan.Agents, Agent{
				Name:        AgentBackendAdversary,
				FileName:    "backend-adversary.agent.md",
				Description: "Challenges implementations and identifies weaknesses.",
			})
		}

		if profile.IsReact {
			plan.Agents = append(plan.Agents, Agent{
				Name:        AgentFrontendAdversary,
				FileName:    "frontend-adversary.agent.md",
				Description: "Challenges implementations and identifies weaknesses.",
			})
		}
	}

	return plan
}

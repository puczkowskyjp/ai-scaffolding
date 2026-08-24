package planning

import "github.com/puczkowskyjp/ai-scaffold/internal/detection"

const (
	AgentArchitect         = "architect"
	AgentDotNet            = "dotnet"
	AgentFrontend          = "frontend"
	AgentDatabase          = "database"
	AgentBackendTesting    = "backend-testing"
	AgentFrontendTesting   = "frontend-testing"
	AgentBackendAdversary  = "backend-adversary"
	AgentFrontendAdversary = "frontend-adversary"
	AgentMicroservices     = "microservices"
)

type AgentPlan struct {
	Agents []Agent
}

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

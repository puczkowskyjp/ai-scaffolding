package detection

// Architecture describes the high-level structural style of a project.
type Architecture string

const (
	// ArchitectureUnknown is used when the architecture has not been identified.
	ArchitectureUnknown Architecture = "unknown"
	// ArchitectureMonolith represents a single-deployment application.
	ArchitectureMonolith Architecture = "monolith"
	// ArchitectureModularMonolith represents a monolith with well-defined internal modules.
	ArchitectureModularMonolith Architecture = "modular-monolith"
	// ArchitectureMicroservices represents a distributed system of independently deployable services.
	ArchitectureMicroservices Architecture = "microservices"
	// ArchitectureServerless represents a function-as-a-service or event-driven deployment model.
	ArchitectureServerless Architecture = "serverless"
)

// ProjectProfile holds the detected characteristics of a repository, including
// the technologies in use and the chosen deployment architecture.
type ProjectProfile struct {
	IsDotNet     bool
	IsReact      bool
	IsVite       bool
	IsPostgres   bool
	Architecture Architecture
	Technologies []string
}

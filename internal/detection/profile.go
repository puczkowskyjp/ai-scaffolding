package detection

// Architecture describes the high-level structural style of a project.
type Architecture string

// Supported Architecture values.
const (
	ArchitectureUnknown         Architecture = "unknown"
	ArchitectureMonolith        Architecture = "monolith"
	ArchitectureModularMonolith Architecture = "modular-monolith"
	ArchitectureMicroservices   Architecture = "microservices"
	ArchitectureServerless      Architecture = "serverless"
)

// ProjectProfile holds the detected characteristics of a repository, including
// the technologies in use and the chosen deployment architecture.
type ProjectProfile struct {
	Name         string
	Languages    []string
	Frameworks   []string
	ProjectTypes []string

	IsDotNet        bool
	IsReact         bool
	IsVite          bool
	IsPostgres      bool
	Architecture    Architecture
	Technologies    []string
	Characteristics []string
}

package detection

type Architecture string

const (
	ArchitectureUnknown         Architecture = "unknown"
	ArchitectureMonolith        Architecture = "monolith"
	ArchitectureModularMonolith Architecture = "modular-monolith"
	ArchitectureMicroservices   Architecture = "microservices"
	ArchitectureServerless      Architecture = "serverless"
)

type ProjectProfile struct {
	IsDotNet     bool
	IsReact      bool
	IsVite       bool
	IsPostgres   bool
	Architecture Architecture
	Technologies []string
}

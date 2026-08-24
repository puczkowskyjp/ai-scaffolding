package planning

// Agent describes a single AI agent to be generated for the repository,
// including its logical name, output file name, and a short human-readable
// description of its role.
type Agent struct {
	Name        string
	FileName    string
	Description string
}

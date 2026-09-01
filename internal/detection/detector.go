package detection

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsDotNetProject reports whether root contains a .NET project by searching for
// .csproj or .sln files.
func IsDotNetProject(root string) bool {
	var found bool

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)

		if ext == ".csproj" || ext == ".sln" {
			found = true
		}

		return nil
	})

	return found
}

// IsViteProject reports whether root contains a Vite project by searching for
// a Vite configuration file (vite.config.ts, vite.config.js, etc.).
func IsViteProject(root string) bool {
	var found bool

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		name := info.Name()

		switch name {
		case "vite.config.ts", "vite.config.js", "vite.config.mts", "vite.config.mjs":
			found = true
		}

		return nil
	})

	return found
}

// IsReactProject reports whether root contains a React project by searching for
// .tsx or .jsx source files.
func IsReactProject(root string) bool {
	var found bool

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)

		if ext == ".tsx" || ext == ".jsx" {
			found = true
		}

		return nil
	})

	return found
}

// IsPostgresProject reports whether root uses PostgreSQL by scanning .csproj and
// docker-compose files for references to npgsql or postgres.
func IsPostgresProject(root string) bool {
	var found bool

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)

		if ext != ".csproj" && info.Name() != "docker-compose.yml" && info.Name() != "compose.yml" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if strings.Contains(strings.ToLower(string(content)), "npgsql") ||
			strings.Contains(strings.ToLower(string(content)), "postgres") {
			found = true
		}

		return nil
	})

	return found
}

// Detect scans the repository at root and returns a ProjectProfile describing
// the technologies detected (e.g. .NET, React, Vite, PostgreSQL).
func Detect(root string) (ProjectProfile, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return ProjectProfile{}, fmt.Errorf("resolve target directory: %w", err)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return ProjectProfile{}, fmt.Errorf("inspect target directory: %w", err)
	}

	if !info.IsDir() {
		return ProjectProfile{}, fmt.Errorf("target path is not a directory: %s", absRoot)
	}

	profile := ProjectProfile{
		Name:         filepath.Base(absRoot),
		IsDotNet:     IsDotNetProject(absRoot),
		IsReact:      IsReactProject(absRoot),
		IsVite:       IsViteProject(absRoot),
		IsPostgres:   IsPostgresProject(absRoot),
		Architecture: ArchitectureUnknown,
	}

	if profile.IsDotNet {
		profile.Languages = append(profile.Languages, "C#")
		profile.Frameworks = append(profile.Frameworks, ".NET")
		profile.ProjectTypes = append(profile.ProjectTypes, "Application")
		profile.Characteristics = append(profile.Characteristics, ".NET application")
		profile.Technologies = append(profile.Technologies, ".NET")
	}

	if profile.IsReact {
		profile.Languages = append(profile.Languages, "JavaScript/TypeScript")
		profile.Frameworks = append(profile.Frameworks, "React")
		profile.ProjectTypes = append(profile.ProjectTypes, "Frontend")
		profile.Characteristics = append(profile.Characteristics, "React frontend")
		profile.Technologies = append(profile.Technologies, "React")
	}

	if profile.IsVite {
		profile.Frameworks = append(profile.Frameworks, "Vite")
		profile.Characteristics = append(profile.Characteristics, "Vite tooling")
		profile.Technologies = append(profile.Technologies, "Vite")
	}

	if profile.IsPostgres {
		profile.Characteristics = append(profile.Characteristics, "PostgreSQL")
		profile.Technologies = append(profile.Technologies, "PostgreSQL")
	}

	return profile, nil
}

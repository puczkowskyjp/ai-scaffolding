package detection

import (
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
func Detect(root string) ProjectProfile {
	profile := ProjectProfile{
		IsDotNet:   IsDotNetProject(root),
		IsReact:    IsReactProject(root),
		IsVite:     IsViteProject(root),
		IsPostgres: IsPostgresProject(root),
	}

	if profile.IsDotNet {
		profile.Technologies = append(profile.Technologies, ".NET")
	}

	if profile.IsReact {
		profile.Technologies = append(profile.Technologies, "React")
	}

	if profile.IsVite {
		profile.Technologies = append(profile.Technologies, "Vite")
	}

	if profile.IsPostgres {
		profile.Technologies = append(profile.Technologies, "PostgreSQL")
	}

	return profile
}

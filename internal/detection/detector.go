package detection

import (
	"os"
	"path/filepath"
	"strings"
)

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

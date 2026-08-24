# Project Instructions

This project uses GitHub Copilot with specialized development agents.

## Architecture

The project architecture is {{.Profile.Architecture}}.

## Detected Technologies

{{range .Profile.Technologies}}
- {{.}}
{{end}}

## Available Agents

{{ range .Plan.Agents }}
- {{ .Name }} — {{ .Description }}
{{ end }}
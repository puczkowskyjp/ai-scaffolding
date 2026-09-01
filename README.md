# ai-scaffolding

`ai-scaffold init` analyzes a repository, proposes Copilot agent scaffolding, and generates the selected files under `.github/`.

## Usage

Interactive workflow:

```bash
go run . init /path/to/repository
```

Non-interactive workflow:

```bash
go run . init --non-interactive /path/to/repository
```

Generated scaffolding includes `.github/copilot-instructions.md` plus agent files under `.github/agents/` based on the selected plan.

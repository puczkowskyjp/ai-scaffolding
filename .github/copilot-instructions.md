(# AI-Scaffold — Copilot Instructions)

This repository implements an AI scaffolding CLI (ai-scaffold) written in Go.
The tool's goal is to reduce boilerplate for agentic coding by initializing a
repository with opinionated agent templates, prompts, and CI-friendly checks.

Goals
- Provide `ai-scaffold init` to scan an existing repository, infer technologies,
	and interactively collect the intended architecture and agent preferences.
- Create or populate a `.github/` layout containing `agents/`, `copilot-instructions.md`,
	and other automation scaffolding so teams can run agent-driven workflows.

CLI behavior (high level)
- `ai-scaffold init` — scans the repo (detects languages, `go.mod`, package layout,
	Dockerfiles, frontends, etc.), then prompts the user for intended architecture,
	whether specialized testing agents are required, CI constraints, and other
	policies.
- The tool is idempotent: it should not overwrite user changes without confirmation.
- After confirmation, it will create or update `.github/agents/` and populate
	files from templates (architect, golang, etc.), and write a repo-level
	`copilot-instructions.md` describing project intent for agents.

Interactive prompts the tool should ask
- Target architecture (monolith, microservices, libraries, CLI, serverless).
- Primary languages and frameworks to treat as first-class (Go, Node, Python, etc.).
- Need for specialized testing agents (unit/contract/integration/e2e).
- CI runners, module/version pinning, and test/coverage expectations.

Outputs and acceptance criteria
- `.github/agents/` created with chosen agent templates.
- A populated `copilot-instructions.md` reflecting repository goals and agent roles.
- Minimal CI workflow(s) added or suggested for running formatters, linters, and tests.
- A short `README.md` note describing how to use generated agents and run `ai-scaffold`.

Developer notes for agents
- Architect agent: leads design discussions, asks clarifying questions, and delegates
	implementation tasks to specialists (for example the Golang agent).
- Golang agent: produces idiomatic, tested Go code, follows `gofmt`/`gofumpt`, `go vet`,
	and writes table-driven tests and benchmarks when necessary.

Edge cases & safety
- When detection is ambiguous, prefer asking the user rather than guessing.
- Do not commit secrets or private tokens. Notify user when repo contains
	credentials-looking files and skip them.

Next steps for implementation
- Implement the `ai-scaffold` CLI in `main.go`, with a subcommand `init` that runs
	a repo scanner and an interactive prompt flow.
- Add templates under `templates/agents/` (architect.agent.md, golang.agent.md) and
	wire them into the generator logic.

Example invocation

```bash
./ai-scaffold init
```

This file documents the project's intent for Copilot/agent workflows and should be
kept concise and human-readable so agents and developers can act on it.


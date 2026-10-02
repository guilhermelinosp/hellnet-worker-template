# hellnet-worker-template

> GitHub template for Go Worker Service projects (Kafka consumers, background workers, etc.).

[![pipeline](https://github.com/guilhermelinosp/hellnet-worker-template/actions/workflows/pipeline.yml/badge.svg)](https://github.com/guilhermelinosp/hellnet-worker-template/actions/workflows/pipeline.yml)
[![pr-check](https://github.com/guilhermelinosp/hellnet-worker-template/actions/workflows/pr-check.yml/badge.svg)](https://github.com/guilhermelinosp/hellnet-worker-template/actions/workflows/pr-check.yml)
[![CodeQL](https://github.com/guilhermelinosp/hellnet-worker-template/actions/workflows/codeql.yml/badge.svg)](https://github.com/guilhermelinosp/hellnet-worker-template/actions/workflows/codeql.yml)

## Initialize from this template

After **Use this template**, clone the new repository and run:

```bash
scripts/init-from-template.sh <repo-name> [service-name]   # renames the module, imports and cmd/ (services)
scripts/setup-repo.sh                                      # repo settings, "main" ruleset and CI variable
```

Then create the `HELLNET_ACTIONS_PRIVATE_KEY` secret (the script prints the exact command) and make sure the
`hellnet-actions` GitHub App is installed on the repository.

## Quick start

```bash
go run ./cmd/worker      # after initialising the repository (see above, `cmd/<service>` once renamed)
```

Without `HELLNET_TELEMETRY_ENDPOINT` the worker still runs and prints a `tick` every 5 s.

## Configuration

| Variable | Purpose | Default |
|---|---|---|
| `HELLNET_TELEMETRY_ENDPOINT` | OTLP/HTTP collector URL (telemetry runs without exporting when empty) | *empty* |
| `HELLNET_TELEMETRY_SERVICE` | service name reported by telemetry | module name |

A `.env` file next to the binary is loaded when present (`internal/env`); variables already set in the environment win.

## Architecture

```text
cmd/worker/main.go   telemetry -> signal-aware context -> worker loop -> graceful shutdown
internal/env         optional .env loading and typed environment helpers
```

`main.go` boots telemetry (`telemetry.New`), runs `workerLoop` in a goroutine (a 5 s ticker calling `runJob`, instrumented as a `tick` worker span) and, on `SIGINT`/`SIGTERM`, cancels the context and waits up to 10 s for the loop to stop. Replace `doWork` with your job (Kafka consumer, queue polling, scheduled task).

## Development

```bash
go test -race ./...
go vet ./...
golangci-lint run ./...
```

Install the git hooks once with `lefthook install`: they run formatting, vet, tests (with and without `-race`), build, `go mod tidy`, lint, `govulncheck` and a secrets scan. Commits follow [Conventional Commits](https://www.conventionalcommits.org/).

## CI/CD

| Workflow | Trigger | What it does |
|---|---|---|
| `pr-check` | pull request | shellcheck, merge strategy and Conventional Commits (`merge-check`), Gitleaks, labels and the Go quality gate (module integrity, vet, race tests with coverage, lint, build, dependency review). `pr-gate` aggregates them and is the required check |
| `pipeline` | push to `main` (ignores `.github/**`) or manual | semver guard (blocks an automatic major), immutable tag + GitHub Release, container image |
| `codeql` | nightly or manual | static analysis (CodeQL) |
| `security` | nightly or manual | Gitleaks and Trivy scans |
| `auto-pr` | push to `feat/**` or `fix/**` | opens the pull request automatically |
| `dependabot-actions-auto-merge` | Dependabot pull requests | auto-merges GitHub Actions bumps |

The workflows call reusable workflows from [templates](https://github.com/guilhermelinosp/templates) at `@latest`. Releases need the `HELLNET_ACTIONS_PRIVATE_KEY` secret and the `HELLNET_ACTIONS_CLIENT_ID` variable (set them with `scripts/setup-repo.sh`).

## Contributing and license

See [CONTRIBUTING.md](CONTRIBUTING.md) and [SECURITY.md](SECURITY.md). Licensed under [Apache 2.0](LICENSE).

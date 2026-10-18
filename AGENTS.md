# AGENTS.md

SIIG is a **single-binary Go CLI**. Treat it as a CLI, not a service.

## Hard constraints

- Stay inside this repo. Do not write to `C:\` unless a toolchain binary already lives there and is required to run.
- Tests and temp files: `X:\ALL FOLDER\TESTING FOR EVERYTHING` (use `gocache`, `gotmp`, `siig-tests`).
- Package / module caches: `X:\ALL FOLDER\MANAGER of all packages and etc`.
- **No Docker. No Kubernetes. No extra infra.** No Dockerfile, compose, helm, or container tooling.
- Do not add dependencies for things the stdlib or existing modules (`cobra`, `fatih/color`) already cover.
- Do not add telemetry, config frameworks, or a logging stack. CLI errors go to stderr; cobra `RunE` returns errors.

## Layout

```
scanner/  walk project, detect stacks
matcher/  merge rules, targets, dedupe
writer/   append ignore files + vscode settings
cache/    ~/.siig/cache, 7-day TTL
fetcher/  optional GitHub raw JSON for unknown stacks
rules/    embedded JSON DB
cmd/      cobra entry
output/   stdout
```

## Commands

```
go test ./...
go vet ./...
go build -o siig .
```

CI must run tests (and vet/lint) on every push and pull request. Release may still cross-compile binaries; that is not “infra.”

## Tests

Put `*_test.go` next to the package. Use `t.TempDir()` / env-overridden temp dirs. Do not write into the user’s real home or real projects.

## Score / mining

Ship one behavior change **with its tests** in a small commit. Do not dump unrelated refactors, formatting, and features in the same commit.

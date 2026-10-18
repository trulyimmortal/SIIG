# Score roadmap (45 → 85+)

DataFactor quality is high enough to keep; **buyer fit fails on tests at HEAD and task capacity ≈ 0**. Quality × zero mineable tasks stays near 45. This file is the work list, not a product spec.

Skipped on purpose (CLI, user rule): Docker, Kubernetes, `.env.example`, Dependabot, structured log library.

## Concept

| Lever | Why it moves score | What we do |
|---|---|---|
| Tests at HEAD | Buyers clone and `go test ./...` | Real assertion tests on matcher, writer, scanner, cache, rules, fetcher validation |
| CI runs tests | `ci_runs_tests` was false | `.github/workflows/ci.yml` on push/PR; release job tests before compile |
| Lint | `has_lint_config=false` | `.golangci.yml` + lint job |
| Docs | no CHANGELOG / CONTRIBUTING | add both; README documents test command |
| Architecture / hygiene | os.Exit + no validation | cobra `RunE`, target parse errors, stack ID / cache path guards, fix dead `HasNpmLib` |
| Task capacity | history must contain small feature+test commits | **commit each item below separately** — a single mega-commit does not mine |

Recency was ~83 days. New commits also refresh that signal.

## Status

### Done in working tree (commit separately)

1. **Scanner `HasNpmLib`** — `package.json` with `main` or `exports` now sets the flag so `.npmignore` can apply. Tests in `scanner/scanner_test.go`.
2. **Matcher** — table tests for `Match` targets, `dedupe`, `ParseTargets`. `matcher/matcher_test.go`.
3. **Writer** — ignore append, dry-run, vscode merge, invalid JSON skip via temp dirs. `writer/writer_test.go`.
4. **Cache** — Get/Set/IsCached/stale with mocked home. `cache/cache_test.go`.
5. **Rules + fetcher** — LoadAll/GetGlobal; reject unsafe stack IDs. No network in tests.
6. **CLI errors** — `RunE` returns errors; unknown `--targets` fail closed.
7. **Lint + CI** — `.golangci.yml`, `ci.yml`, tests before release build.
8. **Docs** — CHANGELOG, CONTRIBUTING, README test/vet commands, this file, `AGENTS.md`.

### How to land (required for 85+)

Do **not** squash into one commit. Suggested sequence:

```
fix: detect npm libraries from package.json main/exports
test: add matcher unit tests for targets and dedupe
test: add writer tests using temp directories
test: add scanner coverage for depth and skip dirs
test: add cache Get/Set/stale tests with mocked HOME
fix: validate stack IDs and unknown CLI targets
ci: run go test and lint on push and pull requests
docs: add CHANGELOG, CONTRIBUTING, and test commands
```

Each commit includes the tests that pin that change.

### After push, re-score

Fresh clone should work with:

```
go test ./...
go vet ./...
go build -o siig .
```

If coverage is still thin, add `output` printer tests and a `cmd` dry-run test — not new packages.

### Will not get you to 85 alone

- More README prose
- Docker
- Fake bulk history
- One giant “quality” commit

# SIIG — She Is Ignoring

> Zero-config CLI that auto-scans your project and populates ignore files based on detected stack.

```
$ siig
✓ Detected: Node.js, Next.js, Docker
✓ .gitignore — 12 entries added
✓ .vercelignore — 4 entries added
✓ .dockerignore — 6 entries added
⚠ .vscode/settings.json — skipped (no .vscode dir found)
```

## What It Does

Run `siig` in any project directory. It:
1. Scans your project (max depth 3)
2. Detects the tech stack (Node, Python, Rust, Go, Next.js, Docker, etc.)
3. Auto-populates `.gitignore`, `.vercelignore`, `.dockerignore`, and `.vscode/settings.json`
4. Never duplicates entries — safe to run multiple times

No configuration. No prompts. Under 200ms.

## Install

**Windows (PowerShell):**
```powershell
irm https://get.siig.dev/install.ps1 | iex
```

**macOS / Linux:**
```bash
curl -fsSL https://get.siig.dev/install.sh | sh
```

**npm:**
```bash
npm install -g siig
```

## Usage

```bash
siig                        # Run in current directory
siig --dry-run              # Preview changes without writing
siig --targets git,vercel   # Only populate specific targets
siig --verbose              # Show every rule applied
siig --version              # Print version
siig --help                 # Print help
```

## Supported Stacks (v1)

| Stack | Detected Via |
|-------|-------------|
| Global (always) | — |
| Node.js | `package.json` |
| Python | `requirements.txt`, `pyproject.toml` |
| Rust | `Cargo.toml` |
| Go | `go.mod` |
| Next.js | `next.config.*` |
| React | `vite.config.*` |
| Tauri | `src-tauri/tauri.conf.json` |
| Android | `AndroidManifest.xml`, `build.gradle` |
| Flutter | `pubspec.yaml` |
| Firebase | `firebase.json` |
| Docker | `Dockerfile`, `docker-compose.yml` |
| Java | `pom.xml`, `build.gradle` |
| .NET / C# | `*.csproj`, `*.sln` |
| Ruby | `Gemfile` |
| PHP | `composer.json` |
| Kotlin | `build.gradle.kts` |
| VS Code Extension | `.vscodeignore` |
| Supabase | `supabase/` |
| Vercel | `vercel.json` |

## Output Targets

| Target | Created When |
|--------|-------------|
| `.gitignore` | Always |
| `.vercelignore` | Vercel or Next.js detected |
| `.dockerignore` | Dockerfile or docker-compose detected |
| `.npmignore` | Library package detected |
| `.vscode/settings.json` | `.vscode/` directory exists |

## Security

- **Fully offline** — no network requests after install
- **No telemetry** — zero data collection
- **Read-only scan** — only reads file/dir names, never file contents
- **Write-safe** — only appends to ignore files, never touches source code
- **Idempotent** — running twice produces the same result
- **Checksum verified** — install scripts verify SHA256 before installing

## Uninstall

**Windows:**
```powershell
irm https://get.siig.dev/uninstall.ps1 | iex
```

**macOS / Linux:**
```bash
curl -fsSL https://get.siig.dev/uninstall.sh | sh
```

## Build from Source

```bash
go build -o siig .
```

## License

MIT

---

Made by [@yv3000](https://github.com/yv3000)

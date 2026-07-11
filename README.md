<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/license-MIT-green?style=for-the-badge" alt="MIT License" />
  <img src="https://img.shields.io/github/v/release/yv3000/SIIG?style=for-the-badge&color=blueviolet" alt="Latest Release" />
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white" alt="Windows" />
  <img src="https://img.shields.io/badge/macOS-000000?style=for-the-badge&logo=apple&logoColor=white" alt="macOS" />
  <img src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black" alt="Linux" />
</p>

<h1 align="center">SIIG — She Is Ignoring</h1>
<p align="center"><b>Zero-config CLI that auto-populates ignore files based on your project's tech stack.</b></p>
<p align="center">One command. No prompts. No config. Under 200ms. 107 stacks supported out of the box.</p>

---

## What is SIIG?

**SIIG** scans your project directory, detects the tech stack, and automatically writes the right entries to `.gitignore`, `.vercelignore`, `.dockerignore`, and `.vscode/settings.json` — so you never accidentally push a `.env` file, API keys, build artifacts, or `node_modules` again.

```
$ siig
✓ Detected: Node.js, Next.js, Docker
✓ .gitignore    — 24 entries added
✓ .vercelignore — 4 entries added
✓ .dockerignore — 8 entries added
⚠ .vscode/settings.json — skipped (no .vscode dir found)
```

- You run `siig` once. It figures out everything else.
- Works on any project — new or existing, any language, any framework.
- **Fully offline** — no network calls, no telemetry, no config files.
- **Idempotent** — running it twice produces the same result. Safe to run anytime.
- **Never touches source code** — only reads file/dir names, only writes to ignore files.

> **Think of it as a one-time setup that actually knows what to ignore.** Instead of copy-pasting from gitignore.io or forgetting half the entries, SIIG just handles it — correctly, instantly, for whatever stack you're running.

---

## Install

SIIG is a single binary. No runtime required. No Node.js, no Python — just download and run.

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/yv3000/SIIG/main/scripts/install.ps1 | iex
```

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/yv3000/SIIG/main/scripts/install.sh | sh
```

The installer will:
1. Detect your OS and architecture (amd64 / arm64)
2. Download the correct binary from GitHub Releases
3. Verify the SHA256 checksum before installing
4. Place `siig` in `~/.siig/bin/` (or `%USERPROFILE%\.siig\bin\` on Windows)
5. Add it to your PATH automatically

---

## Uninstall

### Windows

```powershell
irm https://raw.githubusercontent.com/yv3000/SIIG/main/scripts/uninstall.ps1 | iex
```

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/yv3000/SIIG/main/scripts/uninstall.sh | sh
```

Removes the binary and cleans up the PATH entry. Nothing else is touched — SIIG never modifies anything outside your project directory during normal use.

---

## Proof it works

```bash
# Mixed stack project — Node + Python + Docker + Terraform
$ siig --verbose

✓ Detected: Node.js, Python, Docker, Terraform

Entries by Target:
  gitignore: .env, .env.*, node_modules/, __pycache__/, .terraform/, *.tfstate ...
  dockerignore: node_modules/, __pycache__/, .env, *.pyc ...

✓ .gitignore    — 74 entries added
✓ .dockerignore — 11 entries added
⚠ .vercelignore — skipped (not applicable)
⚠ .npmignore    — skipped (not applicable)

# Run it again — nothing changes
$ siig
✓ Detected: Node.js, Python, Docker, Terraform
✓ Everything already up to date — no changes needed
```

---

## Usage

```bash
siig                          # Auto-detect stack and populate ignore files
siig --dry-run                # Preview what would be added — nothing written
siig --verbose                # Show every rule applied, every entry added
siig --targets git,docker     # Only populate specific targets
siig --version                # Print version
siig --help                   # Print help
```

### Targets

| Flag value | File written |
|---|---|
| `git` | `.gitignore` |
| `vercel` | `.vercelignore` |
| `docker` | `.dockerignore` |
| `npm` | `.npmignore` |
| `vscode` | `.vscode/settings.json` |

Example: `siig --targets git,vercel,docker`

---

## Output Targets

| File | Written When |
|---|---|
| `.gitignore` | Always |
| `.vercelignore` | `vercel.json` or Next.js detected |
| `.dockerignore` | `Dockerfile` or `docker-compose.yml` detected |
| `.npmignore` | Library `package.json` detected (`"main"` or `"exports"` field) |
| `.vscode/settings.json` | `.vscode/` directory exists in project |

---

## Supported Stacks (107)

### Languages

| Stack | Detected Via |
|---|---|
| C | `*.c`, `Makefile`, `CMakeLists.txt` |
| C++ | `*.cpp`, `CMakeLists.txt`, `vcpkg.json` |
| C# | `*.csproj`, `*.sln` |
| .NET | `*.csproj`, `*.sln` |
| F# | `*.fsproj`, `*.fs` |
| Go | `go.mod` |
| Haskell | `*.cabal`, `stack.yaml` |
| Java | `pom.xml`, `build.gradle` |
| Julia | `Project.toml`, `*.jl` |
| Kotlin (JVM) | `build.gradle.kts`, `settings.gradle.kts` |
| Lua | `*.lua`, `*.rockspec` |
| MATLAB | `*.m`, `*.mlx` |
| Nix | `flake.nix`, `default.nix` |
| Objective-C | `*.m`, `Podfile`, `*.xcodeproj` |
| Perl | `*.pl`, `cpanfile` |
| PHP | `composer.json` |
| Python | `requirements.txt`, `pyproject.toml`, `Pipfile` |
| R | `DESCRIPTION`, `*.Rproj`, `renv.lock` |
| Ruby | `Gemfile` |
| Rust | `Cargo.toml` |
| Scala | `build.sbt`, `*.scala` |
| Swift | `Package.swift`, `*.xcodeproj` |
| Zig | `build.zig` |
| Clojure | `project.clj`, `deps.edn` |
| Elixir | `mix.exs` |
| Erlang | `rebar.config`, `*.erl` |
| Groovy | `*.groovy`, `Jenkinsfile` |

### Frontend Frameworks

| Stack | Detected Via |
|---|---|
| Angular | `angular.json` |
| Astro | `astro.config.*` |
| Gatsby | `gatsby-config.*` |
| Next.js | `next.config.*` |
| Nuxt.js | `nuxt.config.*` |
| Qwik | `vite.config.ts` + Qwik |
| React | `vite.config.*` + React |
| Remix | `remix.config.*` |
| SolidJS | `solid.config.*` |
| Svelte / SvelteKit | `svelte.config.*`, `*.svelte` |
| Vue.js | `vue.config.*`, `*.vue` |
| Storybook | `.storybook/main.*` |

### Backend Frameworks

| Stack | Detected Via |
|---|---|
| Django | `manage.py`, `wsgi.py` |
| Express.js | `server.js`, `app.js` |
| FastAPI | `main.py`, `app.py` |
| Flask | `wsgi.py`, `instance/` |
| Laravel | `artisan`, `composer.json` |
| NestJS | `nest-cli.json` |
| Spring Boot | `pom.xml`, `application.properties` |
| Symfony | `symfony.lock` |

### Mobile

| Stack | Detected Via |
|---|---|
| Android | `AndroidManifest.xml`, `build.gradle` |
| Capacitor (Ionic) | `capacitor.config.*` |
| Expo | `app.json`, `eas.json` |
| Flutter | `pubspec.yaml` |
| iOS / Xcode | `*.xcodeproj`, `Info.plist` |
| React Native | `metro.config.*` |
| Tauri | `src-tauri/tauri.conf.json` |

### DevOps & Cloud

| Stack | Detected Via |
|---|---|
| Ansible | `ansible.cfg`, `site.yml` |
| Apache | `.htaccess`, `httpd.conf` |
| AWS CDK | `cdk.json` |
| Docker | `Dockerfile`, `docker-compose.yml` |
| Fly.io | `fly.toml` |
| GitHub Actions | `.github/workflows/` |
| Kubernetes / Helm | `Chart.yaml`, `kustomization.yaml` |
| Netlify | `netlify.toml` |
| Nginx | `nginx.conf` |
| Packer | `*.pkr.hcl` |
| Pulumi | `Pulumi.yaml` |
| Railway | `railway.toml` |
| Terraform | `*.tf`, `*.tfvars` |
| Vagrant | `Vagrantfile` |
| Vercel | `vercel.json` |

### Databases & ORMs

| Stack | Detected Via |
|---|---|
| Drizzle ORM | `drizzle.config.*` |
| MongoDB | `mongod.conf`, `.mongorc.js` |
| MySQL | `my.cnf`, `*.sql` |
| PostgreSQL | `postgresql.conf`, `schema.sql` |
| Prisma | `schema.prisma` |
| Redis | `redis.conf`, `dump.rdb` |
| SQLite | `*.db`, `*.sqlite` |
| Supabase | `supabase/` |

### Game Dev

| Stack | Detected Via |
|---|---|
| Godot | `project.godot`, `*.gd` |
| Unity | `ProjectSettings/`, `Assets/` |
| Unreal Engine | `*.uproject` |

### Build Tools & Monorepos

| Stack | Detected Via |
|---|---|
| Bazel | `WORKSPACE`, `BUILD.bazel` |
| Gradle | `build.gradle`, `gradlew` |
| Lerna | `lerna.json` |
| Maven | `pom.xml`, `.mvn/` |
| Nx | `nx.json` |
| Turborepo | `turbo.json` |

### Runtimes & Package Managers

| Stack | Detected Via |
|---|---|
| Bun | `bun.lockb`, `bunfig.toml` |
| Conda | `environment.yml` |
| Deno | `deno.json`, `deno.lock` |
| Node.js | `package.json` |
| pnpm | `pnpm-workspace.yaml` |
| Yarn | `.yarnrc.yml`, `yarn.lock` |

### Other Tools & Frameworks

| Stack | Detected Via |
|---|---|
| Dart | `pubspec.yaml`, `*.dart` |
| Electron | `electron-builder.yml` |
| Firebase | `firebase.json` |
| GraphQL | `schema.graphql`, `codegen.yml` |
| Jupyter / Data Science | `*.ipynb`, `environment.yml` |
| Kotlin (Android) | `build.gradle`, `AndroidManifest.xml` |
| NestJS | `nest-cli.json` |
| PyTorch / ML | `train.py`, `checkpoints/` |
| Shopify | `shopify.app.toml`, `theme.liquid` |
| Stripe | `stripe.config.*` |
| T3 Stack | `prisma/schema.prisma` + tRPC |
| Tamagui | `tamagui.config.*` |
| Testing (Jest/Vitest/Cypress) | `jest.config.*`, `playwright.config.*` |
| tRPC | `trpc.ts`, `router.ts` |
| Vite | `vite.config.*` |
| VS Code Extension | `.vscodeignore` |
| Wails (Go Desktop) | `wails.json` |
| WordPress | `wp-config.php` |

---

## How It Works

### Architecture

```
siig (single binary, ~5MB)
├── rules/db/       ← 107 JSON rule files, embedded at compile time
├── scanner/        ← walks project dir (max depth 3), collects file/dir names
├── matcher/        ← matches detected files against rule DB, merges entries
└── writer/         ← deduplicates against existing ignore files, appends only new entries
```

### The Flow

```
User runs: siig
      |
      v
  scanner  →  walks project root (max depth 3, skips node_modules/.git)
      |           collects all file names and dir names found
      v
  matcher  →  for each of 107 stack rules:
      |           checks if detect.files or detect.dirs exist in scan result
      |           if match → apply that stack's gitignore/dockerignore/etc rules
      |           global rules always applied regardless of stack
      v
  writer   →  reads existing ignore files
      |           builds a Set of already-present entries
      |           appends only the entries not already there
      |           adds "# Added by SIIG vX.X.X" comment block
      v
  output   →  prints summary of what was detected and written
```

### Rules DB

Each stack is a single JSON file in `rules/db/`. Example — `node.json`:

```json
{
  "id": "node",
  "name": "Node.js",
  "detect": {
    "files": ["package.json", "package-lock.json", "yarn.lock"],
    "dirs": ["node_modules"]
  },
  "rules": {
    "gitignore": ["node_modules/", ".env", ".env.local", "dist/", "*.log"],
    "vercelignore": ["node_modules/", ".env"],
    "dockerignore": ["node_modules/", ".env", "*.log"],
    "vscode": { "files.exclude": { "node_modules": true } }
  }
}
```

All 107 JSON files are embedded into the binary at compile time using Go's `//go:embed` — **no external files needed at runtime.**

---

## Security

- **Fully offline** — zero network requests after install. No API calls, no update checks, nothing.
- **No telemetry** — no data collection of any kind.
- **Read-only scan** — SIIG only reads file and directory *names*, never file contents. Your source code, secrets, and `.env` values are never read.
- **Write-safe** — only appends to ignore files. Never modifies source code, never deletes anything.
- **Idempotent** — every entry is checked against existing content before writing. Running twice is identical to running once.
- **Checksum verified install** — install scripts verify SHA256 of the downloaded binary before placing it on your system.
- **No elevated permissions** — runs entirely as the current user. No admin/sudo required after install.

---

## Platform Support

| Platform | Architecture | Status |
|---|---|---|
| Windows 10/11 | x64 (amd64) | ✅ Fully supported |
| Windows 11 | ARM64 | ✅ Fully supported |
| macOS (Intel) | x64 (amd64) | ✅ Fully supported |
| macOS (Apple Silicon) | ARM64 (arm64) | ✅ Fully supported |
| Linux | x64 (amd64) | ✅ Fully supported |
| Linux | ARM64 | ✅ Fully supported |

---

## Build from Source

```bash
git clone https://github.com/yv3000/SIIG.git
cd SIIG
go build -o siig .
```

Cross-compile for any platform:

```bash
GOOS=windows GOARCH=amd64 go build -o siig-windows-amd64.exe .
GOOS=darwin  GOARCH=arm64 go build -o siig-darwin-arm64 .
GOOS=linux   GOARCH=amd64 go build -o siig-linux-amd64 .
```

**Requirements:** Go 1.21+

---

## FAQ

**Will SIIG overwrite my existing `.gitignore`?**
No. SIIG only appends entries that aren't already present. Your existing content is never modified or removed.

**Is it safe to run on an existing project?**
Yes. SIIG never touches source code — only ignore files. Run `siig --dry-run` first if you want to preview exactly what would be added.

**What if my stack isn't detected?**
Add the detection file SIIG looks for (e.g. `package.json` for Node.js, `Cargo.toml` for Rust). Or run `siig --targets git` to apply global rules only. Community contributions for new stacks are welcome — just add a JSON file to `rules/db/`.

**Does SIIG read my `.env` files or source code?**
No. SIIG only reads file and directory *names* during scanning — never file contents.

**Can I run SIIG in CI?**
Yes — `siig` exits with code 0 on success, non-zero on error. Use `siig --dry-run` in CI to verify ignore files are up to date without writing.

**How do I add a custom stack?**
Open a PR adding a JSON file to `rules/db/` following the schema above.

---

## License

MIT License

Copyright (c) 2026 THE YV

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

---

<p align="center">
  <sub>YV 🖤 ~ I EXPECT NOTHING FROM YOU...</sub>
</p>
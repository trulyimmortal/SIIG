package scanner

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/yv3000/siig/rules"
)

const maxDepth = 3

// DetectedStack holds info about a detected stack
type DetectedStack struct {
	Rule      rules.StackRule
	MatchedBy string // what triggered detection (e.g. "package.json")
}

// ScanResult holds everything the scanner found
type ScanResult struct {
	DetectedStacks []DetectedStack
	HasVSCodeDir   bool
	HasVercelJSON  bool
	HasDockerfile  bool
	HasNpmLib      bool // package.json with "main" or "exports" — it's a library
}

// Scan walks the project directory (max depth 3), collects file/dir names,
// and matches them against detection rules. Caveman approach: just walk and check.
func Scan(projectRoot string, stackRules []rules.StackRule) (*ScanResult, error) {
	if strings.TrimSpace(projectRoot) == "" {
		return nil, fmt.Errorf("project root is required")
	}
	info, err := os.Stat(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("project root: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("project root is not a directory: %s", projectRoot)
	}

	// collect all file names and dir names we find
	fileNames := make(map[string]bool)
	dirNames := make(map[string]bool)
	// also track relative paths for nested detection (e.g. src-tauri/tauri.conf.json)
	relPaths := make(map[string]bool)

	err = filepath.WalkDir(projectRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip errors, keep going caveman style
		}

		// calculate depth relative to project root
		rel, relErr := filepath.Rel(projectRoot, path)
		if relErr != nil {
			return nil
		}

		// skip hidden dirs (except .vscode, .firebase, .dart_tool etc that we detect)
		depth := countSeparators(rel)
		if depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// skip node_modules innards — we only need to know it exists
		if d.IsDir() && d.Name() == "node_modules" {
			dirNames["node_modules"] = true
			return filepath.SkipDir
		}

		// skip .git directory entirely
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		if d.IsDir() {
			dirNames[d.Name()] = true
			// also store relative path for nested dir detection
			relNorm := filepath.ToSlash(rel)
			relPaths[relNorm] = true
		} else {
			fileNames[d.Name()] = true
			// store relative path for nested file detection (e.g. src-tauri/tauri.conf.json)
			relNorm := filepath.ToSlash(rel)
			relPaths[relNorm] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result := &ScanResult{}

	// check for special targets
	result.HasVSCodeDir = dirNames[".vscode"]
	result.HasVercelJSON = fileNames["vercel.json"] || fileNames[".vercelignore"]
	result.HasDockerfile = fileNames["Dockerfile"] || fileNames["docker-compose.yml"] || fileNames["docker-compose.yaml"]
	result.HasNpmLib = detectNpmLib(projectRoot)

	// detect stacks — caveman match: for each rule, check if any detect files/dirs exist
	detectable := rules.GetDetectable(stackRules)
	for _, rule := range detectable {
		if rule.Detect == nil {
			continue
		}

		matched, matchedBy := matchRule(rule, fileNames, dirNames, relPaths)
		if matched {
			result.DetectedStacks = append(result.DetectedStacks, DetectedStack{
				Rule:      rule,
				MatchedBy: matchedBy,
			})
		}
	}

	return result, nil
}

// matchRule checks if any detection file or dir exists in our collected sets
func matchRule(rule rules.StackRule, fileNames, dirNames, relPaths map[string]bool) (bool, string) {
	for _, f := range rule.Detect.Files {
		// check both as bare filename and as relative path
		if fileNames[f] || relPaths[f] {
			return true, f
		}
		// handle glob patterns like *.csproj — just check suffix
		if strings.HasPrefix(f, "*") {
			suffix := f[1:] // e.g. ".csproj"
			for name := range fileNames {
				if strings.HasSuffix(name, suffix) {
					return true, name
				}
			}
		}
	}

	for _, d := range rule.Detect.Dirs {
		if dirNames[d] || relPaths[d] {
			return true, d + "/"
		}
	}

	return false, ""
}

// detectNpmLib is true when root package.json looks like a published library.
func detectNpmLib(projectRoot string) bool {
	data, err := os.ReadFile(filepath.Join(projectRoot, "package.json"))
	if err != nil {
		return false
	}
	var pkg struct {
		Main    string          `json:"main"`
		Exports json.RawMessage `json:"exports"`
	}
	if json.Unmarshal(data, &pkg) != nil {
		return false
	}
	return pkg.Main != "" || len(pkg.Exports) > 0
}

// countSeparators counts path separators to determine depth
func countSeparators(path string) int {
	count := 0
	for _, c := range path {
		if c == filepath.Separator || c == '/' {
			count++
		}
	}
	return count
}

package matcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yv3000/siig/cache"
	"github.com/yv3000/siig/fetcher"
	"github.com/yv3000/siig/rules"
	"github.com/yv3000/siig/scanner"
)

// Target represents one output file target (gitignore, vercelignore, etc.)
type Target string

const (
	TargetGitignore    Target = "gitignore"
	TargetVercelignore Target = "vercelignore"
	TargetDockerignore Target = "dockerignore"
	TargetNpmignore    Target = "npmignore"
	TargetVSCode       Target = "vscode"
)

// AllTargets is the full list of supported targets
var AllTargets = []Target{
	TargetGitignore,
	TargetVercelignore,
	TargetDockerignore,
	TargetNpmignore,
	TargetVSCode,
}

// MatchResult holds the consolidated entries for each target
type MatchResult struct {
	// Entries per target — deduplicated across all detected stacks
	Entries map[Target][]string
	// VSCode settings merged from all stacks
	VSCodeExclude map[string]bool
	// Which targets are applicable based on scan result
	ApplicableTargets map[Target]bool
	// Names of detected stacks for display
	DetectedNames []string
}

// Match takes scan results and all rules, produces a consolidated MatchResult.
// Caveman logic: just loop through everything, merge it all together.
// If offline is false, also tries fetching unknown stacks from siig-rules.
func Match(scanResult *scanner.ScanResult, allRules []StackRule, allowedTargets []Target, offline bool) *MatchResult {
	result := &MatchResult{
		Entries:           make(map[Target][]string),
		VSCodeExclude:     make(map[string]bool),
		ApplicableTargets: make(map[Target]bool),
	}

	// figure out which targets are applicable
	result.ApplicableTargets[TargetGitignore] = true // always
	result.ApplicableTargets[TargetVercelignore] = scanResult.HasVercelJSON || hasStack(scanResult, "nextjs") || hasStack(scanResult, "vercel")
	result.ApplicableTargets[TargetDockerignore] = scanResult.HasDockerfile || hasStack(scanResult, "docker")
	result.ApplicableTargets[TargetVSCode] = scanResult.HasVSCodeDir
	result.ApplicableTargets[TargetNpmignore] = scanResult.HasNpmLib

	// if allowedTargets specified, filter down
	if len(allowedTargets) > 0 {
		filtered := make(map[Target]bool)
		for _, t := range allowedTargets {
			if result.ApplicableTargets[t] {
				filtered[t] = true
			}
		}
		result.ApplicableTargets = filtered
	}

	// collect entries: start with global, then each detected stack
	globalRule := rules.GetGlobal(convertRules(allRules))
	if globalRule != nil {
		mergeRuleEntries(result, globalRule)
	}

	for _, detected := range scanResult.DetectedStacks {
		result.DetectedNames = append(result.DetectedNames, detected.Rule.Name)
		mergeRuleEntries(result, &detected.Rule)
	}

	// online fetch: try to find rules for unknown stacks from siig-rules
	if !offline {
		candidates := buildCandidates(scanResult, allRules)
		if len(candidates) > 0 {
			fetchAndMerge(result, candidates)
		}
	}

	// deduplicate entries per target
	for target := range result.Entries {
		result.Entries[target] = dedupe(result.Entries[target])
	}

	return result
}

// StackRule re-export for convenience (avoiding circular import)
type StackRule = rules.StackRule

// mergeRuleEntries adds all entries from a rule into the result
func mergeRuleEntries(result *MatchResult, rule *rules.StackRule) {
	if len(rule.Rules.Gitignore) > 0 {
		result.Entries[TargetGitignore] = append(result.Entries[TargetGitignore], rule.Rules.Gitignore...)
	}
	if len(rule.Rules.Vercelignore) > 0 {
		result.Entries[TargetVercelignore] = append(result.Entries[TargetVercelignore], rule.Rules.Vercelignore...)
	}
	if len(rule.Rules.Dockerignore) > 0 {
		result.Entries[TargetDockerignore] = append(result.Entries[TargetDockerignore], rule.Rules.Dockerignore...)
	}
	if len(rule.Rules.Npmignore) > 0 {
		result.Entries[TargetNpmignore] = append(result.Entries[TargetNpmignore], rule.Rules.Npmignore...)
	}
	if rule.Rules.VSCode != nil {
		for k, v := range rule.Rules.VSCode.FilesExclude {
			result.VSCodeExclude[k] = v
		}
	}
}

// convertRules converts []StackRule to pass to rules.GetGlobal
// (they're the same type, this is just for the function signature)
func convertRules(allRules []StackRule) []rules.StackRule {
	return allRules
}

// hasStack checks if a stack ID was detected
func hasStack(scanResult *scanner.ScanResult, stackID string) bool {
	for _, s := range scanResult.DetectedStacks {
		if s.Rule.ID == stackID {
			return true
		}
	}
	return false
}

// dedupe removes duplicate strings while preserving order
func dedupe(items []string) []string {
	seen := make(map[string]bool, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// ParseTargets converts user-provided target strings to Target values.
func ParseTargets(targets []string) ([]Target, error) {
	var result []Target
	for _, raw := range targets {
		t := strings.TrimSpace(strings.ToLower(raw))
		if t == "" {
			continue
		}
		switch t {
		case "git", "gitignore":
			result = append(result, TargetGitignore)
		case "vercel", "vercelignore":
			result = append(result, TargetVercelignore)
		case "docker", "dockerignore":
			result = append(result, TargetDockerignore)
		case "npm", "npmignore":
			result = append(result, TargetNpmignore)
		case "vscode":
			result = append(result, TargetVSCode)
		default:
			return nil, fmt.Errorf("unknown target %q (want git, vercel, docker, npm, vscode)", raw)
		}
	}
	return result, nil
}

// genericNames are file/dir names too common to be stack identifiers
var genericNames = map[string]bool{
	"makefile": true, "readme": true, "license": true, "dockerfile": true,
	"index": true, "main": true, "app": true, "src": true, "package": true,
	"config": true, "test": true, "build": true, "dist": true, "lib": true,
	"bin": true, "docs": true, "doc": true, "tmp": true, "temp": true,
	"vendor": true, "public": true, "static": true, "assets": true,
	"changelog": true, "contributing": true, "gitignore": true,
}

// buildCandidates extracts probable stack IDs from scan results that didn't
// match any embedded rule. Caveman: take filename without extension, lowercase it.
func buildCandidates(scanResult *scanner.ScanResult, allRules []StackRule) []string {
	// collect all detect files/dirs from embedded rules to know what's "matched"
	matchedFiles := make(map[string]bool)
	matchedDirs := make(map[string]bool)
	for _, rule := range allRules {
		if rule.Detect == nil {
			continue
		}
		for _, f := range rule.Detect.Files {
			matchedFiles[strings.ToLower(f)] = true
		}
		for _, d := range rule.Detect.Dirs {
			matchedDirs[strings.ToLower(d)] = true
		}
	}

	// also collect IDs of already-detected stacks to skip
	detectedIDs := make(map[string]bool)
	for _, s := range scanResult.DetectedStacks {
		detectedIDs[s.Rule.ID] = true
	}

	// scan all files in the project root for candidate IDs
	// we re-read just the top-level to get candidate names
	candidates := make(map[string]bool)

	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	entries, err := os.ReadDir(cwd)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		name := entry.Name()
		nameLower := strings.ToLower(name)

		// skip hidden files/dirs
		if strings.HasPrefix(name, ".") {
			continue
		}

		// extract candidate ID: name without extension, lowercased
		candidate := strings.TrimSuffix(nameLower, filepath.Ext(nameLower))

		// skip if it's a generic name
		if genericNames[candidate] {
			continue
		}

		// skip if the file is already matched by embedded rules
		if matchedFiles[nameLower] {
			continue
		}
		if entry.IsDir() && matchedDirs[nameLower] {
			continue
		}

		// skip if this ID is already detected
		if detectedIDs[candidate] {
			continue
		}

		// skip very short candidates (1-2 chars are noise)
		if len(candidate) < 3 {
			continue
		}

		candidates[candidate] = true
	}

	// convert to slice
	result := make([]string, 0, len(candidates))
	for c := range candidates {
		result = append(result, c)
	}
	return result
}

// fetchAndMerge fetches rules for candidate stack IDs concurrently and merges them
// into the result. Max 5 concurrent fetches.
func fetchAndMerge(result *MatchResult, candidates []string) {
	const maxConcurrent = 5

	type fetchResult struct {
		rule *rules.StackRule
		id   string
	}

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		sem     = make(chan struct{}, maxConcurrent)
		fetched []fetchResult
	)

	for _, candidate := range candidates {
		wg.Add(1)
		sem <- struct{}{} // acquire semaphore slot

		go func(stackID string) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore slot

			// try cache first
			rule, err := cache.Get(stackID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cache error for %s: %v\n", stackID, err)
			}

			if rule != nil {
				mu.Lock()
				fetched = append(fetched, fetchResult{rule: rule, id: stackID})
				mu.Unlock()
				return
			}

			// cache miss — fetch from network
			rule, err = fetcher.FetchStack(stackID)
			if err != nil {
				// network error — silently skip
				return
			}
			if rule == nil {
				// 404 — stack doesn't exist in siig-rules
				return
			}

			// save to cache (best effort)
			if cacheErr := cache.Set(stackID, rule); cacheErr != nil {
				fmt.Fprintf(os.Stderr, "cache write error for %s: %v\n", stackID, cacheErr)
			}

			mu.Lock()
			fetched = append(fetched, fetchResult{rule: rule, id: stackID})
			mu.Unlock()
		}(candidate)
	}

	wg.Wait()

	// merge fetched rules into result
	for _, fr := range fetched {
		displayName := fr.rule.Name + " (online)"
		result.DetectedNames = append(result.DetectedNames, displayName)
		mergeRuleEntries(result, fr.rule)
	}
}

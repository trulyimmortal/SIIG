package matcher

import (
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
func Match(scanResult *scanner.ScanResult, allRules []StackRule, allowedTargets []Target) *MatchResult {
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

// ParseTargets converts user-provided target strings to Target values
func ParseTargets(targets []string) []Target {
	var result []Target
	for _, t := range targets {
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
		}
	}
	return result
}

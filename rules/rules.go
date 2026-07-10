package rules

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
)

//go:embed db/*.json
var rulesFS embed.FS

// DetectConfig tells the scanner what files/dirs indicate this stack
type DetectConfig struct {
	Files []string `json:"files"`
	Dirs  []string `json:"dirs"`
}

// VSCodeSettings holds vscode-specific settings like files.exclude
type VSCodeSettings struct {
	FilesExclude map[string]bool `json:"files.exclude,omitempty"`
}

// StackRule is one stack's complete rule definition loaded from JSON
type StackRule struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Detect *DetectConfig  `json:"detect"` // nil for global
	Rules  StackRulesList `json:"rules"`
}

// StackRulesList holds all the ignore/config entries for a stack
type StackRulesList struct {
	Gitignore    []string       `json:"gitignore,omitempty"`
	Vercelignore []string       `json:"vercelignore,omitempty"`
	Dockerignore []string       `json:"dockerignore,omitempty"`
	Npmignore    []string       `json:"npmignore,omitempty"`
	VSCode       *VSCodeSettings `json:"vscode,omitempty"`
}

// LoadAll reads every JSON file from the embedded rules/db/ directory
// and returns a slice of StackRule. Caveman style — just read em all, no fancy stuff.
func LoadAll() ([]StackRule, error) {
	var rules []StackRule

	entries, err := fs.ReadDir(rulesFS, "db")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded rules dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := rulesFS.ReadFile("db/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read rule file %s: %w", entry.Name(), err)
		}

		var rule StackRule
		if err := json.Unmarshal(data, &rule); err != nil {
			return nil, fmt.Errorf("failed to parse rule file %s: %w", entry.Name(), err)
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// GetGlobal returns the global rule (detect == nil, always applied)
func GetGlobal(rules []StackRule) *StackRule {
	for i := range rules {
		if rules[i].Detect == nil {
			return &rules[i]
		}
	}
	return nil
}

// GetDetectable returns all rules that have detection config (not global)
func GetDetectable(rules []StackRule) []StackRule {
	var detectable []StackRule
	for _, r := range rules {
		if r.Detect != nil {
			detectable = append(detectable, r)
		}
	}
	return detectable
}

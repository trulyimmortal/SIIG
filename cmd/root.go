package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yv3000/siig/matcher"
	"github.com/yv3000/siig/output"
	"github.com/yv3000/siig/rules"
	"github.com/yv3000/siig/scanner"
	"github.com/yv3000/siig/writer"
)

const appVersion = "1.0.0"

var (
	dryRun  bool
	targets string
	verbose bool
)

// rootCmd is the base command — just run "siig" and it does the thing
var rootCmd = &cobra.Command{
	Use:   "siig",
	Short: "She Is Ignoring — auto-populate ignore files based on detected stack",
	Long: `SIIG scans your project directory, detects the tech stack,
and automatically populates .gitignore, .vercelignore, .dockerignore,
and .vscode/settings.json with the correct entries.

Zero config. One command. No more leaked secrets.`,
	Version: appVersion,
	Run:     runSiig,
}

func init() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be added without writing")
	rootCmd.Flags().StringVar(&targets, "targets", "", "Comma-separated targets to populate (git,vercel,docker,npm,vscode)")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Print every file checked and rule applied")
}

// Execute runs the root command — called from main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// runSiig is the main logic. Caveman flow:
// 1. Load rules
// 2. Scan project
// 3. Match rules
// 4. Write files
// 5. Print what happened
func runSiig(cmd *cobra.Command, args []string) {
	// get current directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot determine current directory: %v\n", err)
		os.Exit(1)
	}

	// 1. load all rules from embedded DB
	allRules, err := rules.LoadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load rules: %v\n", err)
		os.Exit(1)
	}

	// 2. scan project directory
	scanResult, err := scanner.Scan(cwd, allRules)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to scan project: %v\n", err)
		os.Exit(1)
	}

	// print detected stacks
	output.PrintDetected(scanResult)

	if len(scanResult.DetectedStacks) == 0 {
		// still apply global rules
		fmt.Println("Applying global rules only...")
	}

	// 3. parse target filters
	var allowedTargets []matcher.Target
	if targets != "" {
		targetList := strings.Split(targets, ",")
		allowedTargets = matcher.ParseTargets(targetList)
	}

	// 4. match detected stacks against rules
	matchResult := matcher.Match(scanResult, allRules, allowedTargets)

	// verbose output
	if verbose {
		output.PrintVerbose(scanResult, matchResult)
	}

	// 5. write to target files
	writeResults := writer.WriteAll(cwd, matchResult, dryRun)

	// 6. print results
	if dryRun {
		output.PrintDryRun(writeResults, matchResult)
	} else {
		// check if anything was actually written
		anyWritten := false
		for _, r := range writeResults {
			if !r.Skipped && r.EntriesAdded > 0 {
				anyWritten = true
				break
			}
		}

		if anyWritten {
			output.PrintResults(writeResults)
		} else {
			output.PrintNoChanges()
		}
	}
}

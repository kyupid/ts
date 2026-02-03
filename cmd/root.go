package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/kyw/ts/internal/tmux"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ts",
	Short: "tmux session manager",
	Long:  "Create and manage tmux sessions with auto-numbering for same directory",
	RunE:  runRoot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	if !tmux.IsInstalled() {
		return fmt.Errorf("tmux not found")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	baseName := generateSessionName(cwd)
	sessions, err := tmux.ListSessions()
	if err != nil {
		return err
	}

	finalName := findNextAvailable(baseName, sessions)

	if err := tmux.NewSession(finalName, cwd); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return tmux.SwitchOrAttach(finalName)
}

func generateSessionName(path string) string {
	// ~/git/csm-dashboard -> git/csm-dashboard
	home, _ := os.UserHomeDir()
	path = strings.TrimPrefix(path, home+"/")

	parts := strings.Split(path, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return parts[len(parts)-1]
}

func findNextAvailable(baseName string, sessions []tmux.Session) string {
	existing := make(map[string]bool)
	for _, s := range sessions {
		existing[s.Name] = true
	}

	if !existing[baseName] {
		return baseName
	}

	// baseName-2, baseName-3, ...
	re := regexp.MustCompile(`^` + regexp.QuoteMeta(baseName) + `(-(\d+))?$`)
	maxNum := 1

	for _, s := range sessions {
		if matches := re.FindStringSubmatch(s.Name); matches != nil {
			if matches[2] != "" {
				if num, _ := strconv.Atoi(matches[2]); num > maxNum {
					maxNum = num
				}
			}
		}
	}

	return fmt.Sprintf("%s-%d", baseName, maxNum+1)
}

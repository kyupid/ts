package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/kyupid/ts/internal/tmux"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "ts [session-name]",
	Short:        "tmux session manager",
	Long:         "Create and manage tmux sessions with auto-numbering for same directory",
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(1),
	RunE:         runRoot,
}

func Execute() {
	// Handle session name argument before Cobra's subcommand matching
	// Cobra treats unknown first args as subcommands, so we intercept here
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if !isKnownCommand(arg) && !strings.HasPrefix(arg, "-") {
			if err := runRoot(rootCmd, []string{arg}); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			return
		}
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func isKnownCommand(arg string) bool {
	knownCmds := []string{"list", "l", "ls", "switch", "s", "kill", "k", "help", "completion"}
	for _, cmd := range knownCmds {
		if arg == cmd {
			return true
		}
	}
	return false
}

func completeSessionNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	sessions, err := tmux.ListSessions()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names := []string{"all"}
	for _, s := range sessions {
		names = append(names, s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func runRoot(cmd *cobra.Command, args []string) error {
	if !tmux.IsInstalled() {
		return fmt.Errorf("tmux not found")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	var finalName string
	if len(args) > 0 {
		// 사용자가 세션 이름을 지정한 경우
		finalName = strings.ReplaceAll(args[0], ".", "_")
	} else {
		// 자동 이름 생성
		baseName := generateSessionName(cwd)
		sessions, err := tmux.ListSessions()
		if err != nil {
			return err
		}
		finalName = findNextAvailable(baseName, sessions)
	}

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
	var name string
	if len(parts) >= 2 {
		name = parts[len(parts)-2] + "/" + parts[len(parts)-1]
	} else {
		name = parts[len(parts)-1]
	}
	// tmux converts . to _ in session names
	return strings.ReplaceAll(name, ".", "_")
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

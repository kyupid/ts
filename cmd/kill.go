package cmd

import (
	"fmt"
	"os/exec"

	"github.com/kyupid/ts/internal/tmux"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:     "kill [session-name]",
	Aliases: []string{"k"},
	Short:   "Kill a session (fzf if no name given)",
	RunE:    runKill,
}

func init() {
	rootCmd.AddCommand(killCmd)
}

func runKill(cmd *cobra.Command, args []string) error {
	if !tmux.IsInstalled() {
		return fmt.Errorf("tmux not found")
	}

	var target string

	if len(args) > 0 {
		target = args[0]
	} else {
		if _, err := exec.LookPath("fzf"); err != nil {
			return fmt.Errorf("fzf not found")
		}

		sessions, err := tmux.ListSessions()
		if err != nil {
			return err
		}

		if len(sessions) == 0 {
			fmt.Println("no sessions")
			return nil
		}

		var names []string
		for _, s := range sessions {
			names = append(names, s.Name)
		}

		selected, err := fzfSelect(names, "kill session")
		if err != nil {
			return nil // user cancelled
		}
		target = selected
	}

	if err := tmux.KillSession(target); err != nil {
		return fmt.Errorf("failed to kill session %s: %w", target, err)
	}

	fmt.Printf("killed: %s\n", target)
	return nil
}

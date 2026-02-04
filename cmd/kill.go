package cmd

import (
	"fmt"

	"github.com/kyupid/ts/internal/tmux"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:     "kill [session-name]",
	Aliases: []string{"k"},
	Short:   "Kill current session (or specified session)",
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
		if !tmux.IsInsideTmux() {
			return fmt.Errorf("not inside tmux session")
		}
		current, err := tmux.CurrentSession()
		if err != nil {
			return fmt.Errorf("failed to get current session: %w", err)
		}
		target = current
	}

	if err := tmux.KillSession(target); err != nil {
		return fmt.Errorf("failed to kill session %s: %w", target, err)
	}

	fmt.Printf("killed: %s\n", target)
	return nil
}

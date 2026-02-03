package cmd

import (
	"fmt"

	"github.com/kyupid/ts/internal/tmux"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tmux sessions",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	if !tmux.IsInstalled() {
		return fmt.Errorf("tmux not found")
	}

	sessions, err := tmux.ListSessions()
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Println("no sessions")
		return nil
	}

	for _, s := range sessions {
		fmt.Println(s.Name)
	}

	return nil
}

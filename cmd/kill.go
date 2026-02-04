package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kyupid/ts/internal/tmux"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:               "kill [session-name]",
	Aliases:           []string{"k"},
	Short:             "Kill current session (or specified session)",
	RunE:              runKill,
	ValidArgsFunction: completeSessionNames,
}

func init() {
	rootCmd.AddCommand(killCmd)
}

func runKill(cmd *cobra.Command, args []string) error {
	if !tmux.IsInstalled() {
		return fmt.Errorf("tmux not found")
	}

	if len(args) > 0 && args[0] == "all" {
		return killAllSessions()
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

func killAllSessions() error {
	sessions, err := tmux.ListSessions()
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Println("no sessions to kill")
		return nil
	}

	fmt.Printf("sessions to kill (%d):\n", len(sessions))
	for _, s := range sessions {
		fmt.Printf("  - %s\n", s.Name)
	}

	fmt.Print("\nkill all sessions? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		fmt.Println("cancelled")
		return nil
	}

	for _, s := range sessions {
		if err := tmux.KillSession(s.Name); err != nil {
			fmt.Printf("failed to kill %s: %v\n", s.Name, err)
			continue
		}
		fmt.Printf("killed: %s\n", s.Name)
	}

	return nil
}

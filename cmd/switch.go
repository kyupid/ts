package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kyupid/ts/internal/tmux"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:     "switch",
	Aliases: []string{"s"},
	Short:   "Switch to a session using fzf",
	RunE:    runSwitch,
}

func init() {
	rootCmd.AddCommand(switchCmd)
}

func runSwitch(cmd *cobra.Command, args []string) error {
	if !tmux.IsInstalled() {
		return fmt.Errorf("tmux not found")
	}

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

	selected, err := fzfSelect(names, "switch to session")
	if err != nil {
		return nil // user cancelled
	}

	return tmux.SwitchOrAttach(selected)
}

func fzfSelect(items []string, header string) (string, error) {
	cmd := exec.Command("fzf", "--reverse", "--header", header)
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

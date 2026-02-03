package tmux

import (
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Session struct {
	Name string
}

func IsInstalled() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

func IsInsideTmux() bool {
	return os.Getenv("TMUX") != ""
}

func ListSessions() ([]Session, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
	out, err := cmd.Output()
	if err != nil {
		// no server running 등은 빈 리스트로 처리
		if strings.Contains(err.Error(), "exit status") {
			return []Session{}, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	sessions := make([]Session, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			sessions = append(sessions, Session{Name: line})
		}
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Name < sessions[j].Name
	})

	return sessions, nil
}

func NewSession(name, path string) error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", name, "-c", path)
	return cmd.Run()
}

func KillSession(name string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", name)
	return cmd.Run()
}

func SwitchClient(name string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", name)
	return cmd.Run()
}

func AttachSession(name string) error {
	cmd := exec.Command("tmux", "attach-session", "-t", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SwitchOrAttach(name string) error {
	if IsInsideTmux() {
		return SwitchClient(name)
	}
	return AttachSession(name)
}

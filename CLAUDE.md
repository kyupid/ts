# ts - tmux session manager

## Philosophy

- **Do one thing well**: Create, switch, kill sessions. Everything else is tmux's job.
- **Zero config**: No config files. Directory-based auto-naming.
- **Tmux native**: Just a wrapper. No daemons, no state.

## Build & Test

```bash
go build -o ts_test    # Build binary
go test ./...          # Run tests
make install           # Install to ~/bin
```

## Project Structure

```
cmd/           # Cobra commands (root, list, switch, kill)
internal/tmux/ # tmux interaction layer
```

## Code Style

- Keep functions small and focused
- Error messages: lowercase, no period (e.g., `fmt.Errorf("tmux not found")`)
- Use `tmux.` package for all tmux operations

## Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `ts` | - | Create new session |
| `ts list` | `l`, `ls` | List sessions |
| `ts switch` | `s` | Switch session (fzf) |
| `ts kill` | `k` | Kill session |

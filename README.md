# ts

tmux session manager with auto-numbering.

## Philosophy

```
Do one thing well    Create, switch, kill sessions. Everything else is tmux's job.
Zero config          No config files. Directory-based auto-naming.
Tmux native          Just a wrapper. No daemons, no state.
```

## Install

```bash
go install github.com/kyupid/ts@latest
```

Or build from source:

```bash
git clone https://github.com/kyupid/ts.git
cd ts
make install  # installs to ~/bin
```

## Usage

```bash
# Create session (auto-named from current directory)
ts                    # ~/git/myproject → git/myproject

# Create session with custom name
ts myproject          # → myproject

# Create another in same directory
ts                    # → git/myproject-2

# List sessions
ts list               # aliases: l, ls

# Switch session (fzf)
ts switch             # alias: s

# Kill session
ts kill               # fzf selection (alias: k)
ts kill git/myproject # by name
```

## Tmux Keybindings

Add to `~/.tmux.conf`:

```bash
bind-key "T" display-popup -E -w 60% -h 60% "ts switch"
bind-key "X" display-popup -E -w 60% -h 60% "ts kill"
```

## Requirements

- tmux
- fzf (for switch/kill)
- Go 1.21+

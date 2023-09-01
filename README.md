# Git Branch TUI

## Description
Switch Git branches interactively via TUI. An improvement for the existing `git branch` command.

## Features
- List Git branches in a specific repo
- Checkout the selected branch with a single keystroke.

## Usage
1. Clone or download the repo
   - `git clone https://github.com/jonasrdl/git-branch-tui`
2. Build the application
   - `go build -o git-branch-tui main.go`

## Creating an alias
If you still want to `git branch` command, you can create an alias to this command.   
Means, if you run `git branch`, its "forwarded" to this application.   

You can do this like following:

Add this line to your e.g. `~/.zshrc` / `~/.bashrc` (depends on which shell you are using)   
`$REPO_PATH` is the path, of this repo on you local machine. (e.g. `/home/user/projects/git-branch-tui`)
```
alias git="git-branch-switcher() { if [[ \"\$1\" == \"branch\" ]]; then $REPO_PATH/git-branch-tui; else command git \"\$@\"; fi; }; git-branch-switcher"
```

Then just restart your terminal emulator, or run `source ~/.zshrc` / `source ~/.bashrc`

Now, you can use the git branch command as you normally would, and it will use this application instead.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// Parse command-line arguments.
	path := flag.String("path", ".", "Path where the Git command should run.")
	flag.Parse()

	// Ensure the specified path exists and is a directory.
	if !isValidDirectory(*path) {
		fmt.Println("Invalid path:", *path)
		return
	}

	app := tview.NewApplication()

	// Create a table to display branches.
	branchTable := tview.NewTable()

	// Fetch Git branches and populate the table.
	branches, err := getGitBranches(*path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for i, branch := range branches {
		branchTable.SetCell(i, 0, tview.NewTableCell(branch).SetAlign(tview.AlignLeft))
	}

	// Initialize the selected row.
	selectedRow := 0
	var selectedBranch string // Variable to store the selected branch.

	// Handle branch checkout and exit the TUI.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			selectedBranch = branchTable.GetCell(selectedRow, 0).Text
			if selectedBranch != "" {
				// Run `git checkout` command to switch to the selected branch.
				cmd := exec.Command("git", "checkout", selectedBranch)
				cmd.Dir = *path
				if err := cmd.Run(); err != nil {
					fmt.Println("Error:", err)
				} else {
					app.Stop() // Exit the TUI
				}
			}
		case tcell.KeyUp:
			if selectedRow > 0 {
				branchTable.GetCell(selectedRow, 0).SetBackgroundColor(tcell.ColorDefault)
				selectedRow--
				branchTable.GetCell(selectedRow, 0).SetBackgroundColor(tcell.ColorYellow)
			}
		case tcell.KeyDown:
			if selectedRow < len(branches)-1 {
				branchTable.GetCell(selectedRow, 0).SetBackgroundColor(tcell.ColorDefault)
				selectedRow++
				branchTable.GetCell(selectedRow, 0).SetBackgroundColor(tcell.ColorYellow)
			}
		}
		return event
	})

	// Set the initial selection.
	branchTable.GetCell(selectedRow, 0).SetBackgroundColor(tcell.ColorYellow)

	// Create a Flex layout to arrange widgets.
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false). // Background box to fill empty space.
		AddItem(branchTable, 0, 1, true)      // Align the table to the bottom.

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	// After the TUI exits, print out the selected branch.
	if selectedBranch != "" {
		fmt.Println("Checked out branch:", selectedBranch)
	}
}

func getGitBranches(path string) ([]string, error) {
	// Run `git branch --list` command to get a list of branches.
	cmd := exec.Command("git", "branch", "--list")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse the output and extract branch names.
	branches := []string{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		branchName := strings.TrimSpace(strings.TrimPrefix(line, "*"))
		if branchName != "" {
			branches = append(branches, branchName)
		}
	}

	return branches, nil
}

func isValidDirectory(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

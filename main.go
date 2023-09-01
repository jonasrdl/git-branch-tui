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
	path := flag.String("path", ".", "Path where the Git command should run.")
	flag.Parse()

	if !isValidDirectory(*path) {
		fmt.Println("Invalid path:", *path)
		return
	}

	app := tview.NewApplication()

	branchTable := tview.NewTable()

	branches, err := getGitBranches(*path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for i, branch := range branches {
		branchTable.SetCell(i, 0, tview.NewTableCell(branch).SetAlign(tview.AlignLeft))
	}

	selectedRow := 0
	var selectedBranch string

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			selectedBranch = branchTable.GetCell(selectedRow, 0).Text
			if selectedBranch != "" {
				cmd := exec.Command("git", "checkout", selectedBranch)
				cmd.Dir = *path
				if err := cmd.Run(); err != nil {
					fmt.Println("Error:", err)
				} else {
					app.Stop()
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

	branchTable.GetCell(selectedRow, 0).SetBackgroundColor(tcell.ColorYellow)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(branchTable, 0, 1, true)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	if selectedBranch != "" {
		fmt.Println("Checked out branch:", selectedBranch)
	}
}

func getGitBranches(path string) ([]string, error) {
	cmd := exec.Command("git", "branch", "--list")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

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

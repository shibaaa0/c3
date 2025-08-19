package main

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FileManager holds the state of the c3 file manager
type FileManager struct {
	app        *tview.Application
	list       *tview.List
	preview    *tview.TextView
	currentDir string
}

// NewFileManager initializes a new file manager
func NewFileManager() *FileManager {
	fm := &FileManager{
		app:        tview.NewApplication(),
		list:       tview.NewList().ShowSecondaryText(false),
		preview:    tview.NewTextView().SetText("Preview pane (to be implemented)"),
		currentDir: getCurrentDir(),
	}
	return fm
}

// getCurrentDir returns the current working directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err) // Simplify for Day 1; we'll handle errors later
	}
	return dir
}

// setupUI configures the TUI layout and starts the app
func (fm *FileManager) setupUI() {
	// Create Flex layout: split screen
	flex := tview.NewFlex().
		AddItem(fm.list, 0, 1, true). // Left: List, dynamic width, focused
		AddItem(fm.preview, 0, 1, false) // Right: Preview, dynamic width

	// Style the list
	fm.list.SetBorder(true).SetTitle("c3 - File Manager").SetTitleAlign(tview.AlignLeft)

	// Populate list with files
	fm.updateFileList()

	// Set key bindings
	fm.setupKeyBindings()

	// Run the app
	if err := fm.app.SetRoot(flex, true).SetFocus(fm.list).Run(); err != nil {
		panic(err)
	}
}

// updateFileList populates the list with files and directories
func (fm *FileManager) updateFileList() {
	fm.list.Clear()
	dirEntries, err := os.ReadDir(fm.currentDir)
	if err != nil {
		fm.list.AddItem("Error reading directory", err.Error(), 0, nil)
		return
	}

	// Sort entries: directories first, then files
	var dirs, files []os.DirEntry
	for _, entry := range dirEntries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	// Add parent directory
	fm.list.AddItem("..", "", 0, fm.navigateToParent)

	// Add directories
	for _, dir := range dirs {
		fm.list.AddItem(dir.Name()+"/", "", 0, func() {
			fm.navigateTo(dir.Name())
		})
	}

	// Add files
	for _, file := range files {
		fm.list.AddItem(file.Name(), "", 0, nil) // No action for files yet
	}

	// Update preview with current directory
	fm.preview.SetText("Current directory: " + fm.currentDir)
}

// navigateTo changes to the specified directory
func (fm *FileManager) navigateTo(dir string) {
	newPath := filepath.Join(fm.currentDir, dir)
	if info, err := os.Stat(newPath); err == nil && info.IsDir() {
		fm.currentDir = newPath
		fm.updateFileList()
		fm.app.SetFocus(fm.list)
	}
}

// navigateToParent navigates to the parent directory
func (fm *FileManager) navigateToParent() {
	fm.navigateTo("..")
}

// setupKeyBindings configures key bindings for navigation
func (fm *FileManager) setupKeyBindings() {
	fm.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			selected := fm.list.GetCurrentItem()
			mainText, _ := fm.list.GetItemText(selected)
			if selected == 0 { // Parent directory
				fm.navigateToParent()
			} else if len(mainText) > 0 && mainText[len(mainText)-1] == '/' {
				fm.navigateTo(mainText[:len(mainText)-1])
			}
			return nil
		case tcell.KeyEscape:
			fm.app.Stop()
			return nil
		}
		return event
	})
}

func main() {
	fm := NewFileManager()
	fm.setupUI()
}
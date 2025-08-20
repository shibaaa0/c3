package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FileManager holds the state of the c3 file manager
type FileManager struct {
	app         *tview.Application
	list        *tview.List
	preview     *tview.TextView
	searchInput *tview.InputField
	listFlex    *tview.Flex
	currentDir  string
	showHidden  bool
	searchQuery string
	searchActive bool
}

// NewFileManager initializes a new file manager
func NewFileManager() *FileManager {
	fm := &FileManager{
		app:         tview.NewApplication(),
		list:        tview.NewList().ShowSecondaryText(false),
		preview:     tview.NewTextView().SetText("Preview pane (to be implemented)").SetDynamicColors(true),
		searchInput: tview.NewInputField().SetLabel("Search: "),
		listFlex:    tview.NewFlex().SetDirection(tview.FlexRow),
		currentDir:  getCurrentDir(),
		showHidden:  false, // Hide hidden files by default
		searchQuery: "",
		searchActive: false,
	}
	// Initially only add list to listFlex
	fm.listFlex.AddItem(fm.list, 0, 1, true)
	return fm
}

// getCurrentDir returns the current working directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err) // Simplify for Day 2; we'll handle errors later
	}
	return dir
}

// setupUI configures the TUI layout and starts the app
func (fm *FileManager) setupUI() {
	// Create main flex: listFlex on left, preview on right
	flex := tview.NewFlex().
		AddItem(fm.listFlex, 0, 1, true).
		AddItem(fm.preview, 0, 1, false)

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
		fm.list.AddItem("[red]Error reading directory[-]", err.Error(), 0, nil)
		return
	}

	// Sort entries: directories first, then files
	var dirs, files []os.DirEntry
	for _, entry := range dirEntries {
		if !fm.showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue // Skip hidden files if showHidden is false
		}
		if fm.searchQuery != "" && !strings.HasPrefix(strings.ToLower(entry.Name()), strings.ToLower(fm.searchQuery)) {
			continue // Skip entries not starting with search query
		}
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	// Add directories
	for _, dir := range dirs {
		fm.list.AddItem("[blue]"+dir.Name()+"/[-]", "", 0, func() {
			fm.navigateTo(dir.Name())
		})
	}

	// Add files
	for _, file := range files {
		fm.list.AddItem("[white]"+file.Name()+"[-]", "", 0, nil) // No action for files yet
	}

	// Update list title with abbreviated path (last 2 directories)
	promptPath := fm.currentDir
	parts := strings.Split(promptPath, string(os.PathSeparator))
	if len(parts) > 2 {
		promptPath = strings.Join(parts[len(parts)-2:], "/")
	} else if strings.HasPrefix(promptPath, "/") {
		promptPath = promptPath[1:] // Remove leading slash for short paths
	}
	fm.list.SetTitle("c3 - " + promptPath)

	// Update preview with current directory, hidden status, and search query
	hiddenStatus := "Hidden files: " + map[bool]string{true: "shown", false: "hidden"}[fm.showHidden]
	searchStatus := "Search: " + map[string]string{"": "none", fm.searchQuery: fm.searchQuery}[fm.searchQuery]
	fm.preview.SetText("Current directory: " + fm.currentDir + "\n" + hiddenStatus + "\n" + searchStatus)
}

// navigateTo changes to the specified directory
func (fm *FileManager) navigateTo(dir string) {
	newPath := filepath.Join(fm.currentDir, dir)
	if info, err := os.Stat(newPath); err == nil && info.IsDir() {
		fm.currentDir = newPath
		fm.searchQuery = ""
		fm.searchActive = false
		fm.listFlex.RemoveItem(fm.searchInput)
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
		// If search is active, only allow specific keys to pass to app
		if fm.searchActive {
			switch event.Key() {
			case tcell.KeyEnter, tcell.KeyEscape:
				// Handle in searchInput's DoneFunc
				return event
			default:
				// Let search input handle other keys
				return event
			}
		}

		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyRight:
			selected := fm.list.GetCurrentItem()
			mainText, _ := fm.list.GetItemText(selected)
			// Remove color tags for navigation
			cleanText := strings.TrimPrefix(strings.TrimSuffix(mainText, "[-]"), "[blue]")
			cleanText = strings.TrimPrefix(cleanText, "[white]")
			if len(cleanText) > 0 && cleanText[len(cleanText)-1] == '/' {
				fm.navigateTo(cleanText[:len(cleanText)-1])
			}
			return nil
		case tcell.KeyLeft:
			fm.navigateToParent()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 's', 'S': // Toggle hidden files with S
				fm.showHidden = !fm.showHidden
				fm.updateFileList()
				return nil
			case 'f', 'F': // Start search with F
				fm.searchActive = true
				fm.listFlex.AddItem(fm.searchInput, 1, 1, false)
				fm.app.SetFocus(fm.searchInput)
				return nil
			}
		case tcell.KeyEscape:
			fm.app.Stop()
			return nil
		}
		return event
	})

	// Search input key bindings
	fm.searchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			fm.searchQuery = strings.TrimSpace(fm.searchInput.GetText())
			fm.searchActive = false
			fm.listFlex.RemoveItem(fm.searchInput)
			fm.updateFileList()
			fm.app.SetFocus(fm.list)
		} else if key == tcell.KeyEscape {
			fm.searchQuery = ""
			fm.searchActive = false
			fm.listFlex.RemoveItem(fm.searchInput)
			fm.searchInput.SetText("")
			fm.updateFileList()
			fm.app.SetFocus(fm.list)
		}
	}).SetChangedFunc(func(text string) {
		fm.searchQuery = strings.TrimSpace(text)
		fm.updateFileList()
	})
}

func main() {
	fm := NewFileManager()
	fm.setupUI()
}
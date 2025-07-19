package main

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	rootDir, _ := os.Getwd()
	root := tview.NewTreeNode(rootDir).SetReference(rootDir).SetExpanded(true)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// Đệ quy thêm node con
	add := func(target *tview.TreeNode, path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			return
		}
		for _, file := range files {
			childPath := filepath.Join(path, file.Name())
			node := tview.NewTreeNode(file.Name()).SetReference(childPath)
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
				node.SetExpanded(false)
			}
			target.AddChild(node)
		}
	}

	// Lần đầu load thư mục root
	add(root, rootDir)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		path := node.GetReference().(string)
		info, err := os.Stat(path)
		if err != nil {
			return
		}
		if info.IsDir() {
			if len(node.GetChildren()) == 0 {
				add(node, path)
			}
			node.SetExpanded(!node.IsExpanded())
		}
	})

	if err := app.SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}

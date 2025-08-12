package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"path/filepath"
)

func main() {
	// Tạo ứng dụng tview
	app := tview.NewApplication()

	// Tạo TreeView để hiển thị tệp/thư mục
	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle("File Manager")

	// Lấy thư mục hiện tại làm thư mục gốc
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root := tview.NewTreeNode(filepath.Base(rootDir)).
		SetReference(rootDir).
		SetColor(tcell.ColorBlue)
	tree.SetRoot(root).SetCurrentNode(root)

	// Hàm để tải danh sách tệp/thư mục
	populateTree := func(node *tview.TreeNode) {
		path := node.GetReference().(string)
		files, err := os.ReadDir(path)
		if err != nil {
			return
		}
		for _, file := range files {
			child := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name()))
			if file.IsDir() {
				child.SetColor(tcell.ColorGreen).
					SetSelectable(true)
			}
			node.AddChild(child)
		}
	}

	// Tải thư mục gốc
	populateTree(root)

	// Xử lý sự kiện khi chọn node
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return
		}
		path := reference.(string)
		fileInfo, err := os.Stat(path)
		if err != nil {
			return
		}
		if fileInfo.IsDir() {
			node.ClearChildren()
			populateTree(node)
			node.SetExpanded(!node.IsExpanded())
		}
	})

	// Hàm hỗ trợ lấy tất cả node
	getAllNodes := func(root *tview.TreeNode) []*tview.TreeNode {
		var nodes []*tview.TreeNode
		var walk func(*tview.TreeNode)
		walk = func(node *tview.TreeNode) {
			nodes = append(nodes, node)
			for _, child := range node.GetChildren() {
				if child.IsExpanded() {
					walk(child)
				}
			}
		}
		walk(root)
		return nodes
	}

	// Hàm hỗ trợ di chuyển xuống
	getNextNode := func(current, root *tview.TreeNode) *tview.TreeNode {
		if current == nil {
			return root
		}
		nodes := getAllNodes(root)
		for i, node := range nodes {
			if node == current && i < len(nodes)-1 {
				return nodes[i+1]
			}
		}
		return current
	}

	// Hàm hỗ trợ di chuyển lên
	getPreviousNode := func(current, root *tview.TreeNode) *tview.TreeNode {
		if current == nil {
			return root
		}
		nodes := getAllNodes(root)
		for i, node := range nodes {
			if node == current && i > 0 {
				return nodes[i-1]
			}
		}
		return current
	}

	// Thiết lập phím tắt
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q': // Thoát ứng dụng
				app.Stop()
				return nil
			case 'j': // Di chuyển xuống
				current := tree.GetCurrentNode()
				if current != nil {
					tree.SetCurrentNode(getNextNode(current, tree.GetRoot()))
				}
				return nil
			case 'k': // Di chuyển lên
				current := tree.GetCurrentNode()
				if current != nil {
					tree.SetCurrentNode(getPreviousNode(current, tree.GetRoot()))
				}
				return nil
			}
		}
		return event
	})

	// Tạo layout chính
	flex := tview.NewFlex().
		AddItem(tree, 0, 1, true)

	// Chạy ứng dụng
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
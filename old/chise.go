package main

import (
	//"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	//"os"
	//"path/filepath"
)

func main() {
	// Tạo ứng dụng tview
	app := tview.NewApplication()

	// Tạo TreeView để hiển thị tệp/thư mục
	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle("File Manager")

	flex := tview.NewFlex().
		AddItem(tree, 0, 1, true)

	// Chạy ứng dụng
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
	
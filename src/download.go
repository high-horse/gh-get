package main

import (
	"os"
	"fmt"
	"log"

	"path/filepath"
	"github.com/rivo/tview"
)


func handleDownload(tree *tview.TreeView) error {
	var selected []*Content

	root := tree.GetRoot()
	if root == nil {
		return fmt.Errorf("tree is empty")
	}

	collectSelectedFromNode(root, &selected)

	if len(selected) == 0 {
		return fmt.Errorf("no files selected")
	}

	baseDir := fmt.Sprintf("%s-%s", reponame, selectedBranch)

	for _, file := range selected {
		if file.IsDir {
			continue
		}

		localPath := filepath.Join(baseDir, file.Path)

		if err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); err != nil {
			return err
		}

		log.Println("Downloading:", file.Path)

		if err := downloadFile(file.DownloadUrl, localPath); err != nil {
			log.Println("failed:", file.Path, err)
			continue
		}
	}

	log.Println("Download complete")
	return nil
}



func collectSelectedFromNode(node *tview.TreeNode, result *[]*Content) {
	ref := node.GetReference()
	if ref != nil {
		c := ref.(*Content)
		if !c.IsDir && c.Selected {
			*result = append(*result, c)
		}
	}

	for _, child := range node.GetChildren() {
		collectSelectedFromNode(child, result)
	}
}
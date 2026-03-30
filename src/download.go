package main

import (
	"fmt"
	// "log"
	"os"
	"path/filepath"
	"sync"

	"github.com/rivo/tview"
)
func handleDownload(tree *tview.TreeView) error  {
	if preserveDirTree {
		return handleDownloadWithDirTreePreserved(tree)
	} 
	return handleDownloadWithoutDirTreePreserved(tree)
}

func handleDownloadWithDirTreePreserved(tree *tview.TreeView) error {
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

	var wg sync.WaitGroup
	var mu sync.Mutex // to safely append to errors slice
	var errors []error

	for _, file := range selected {
		if file.IsDir {
			continue
		}

		wg.Add(1)

		// capture the variable for closure
		file := file

		go func() {
			defer wg.Done()

			localPath := filepath.Join(baseDir, file.Path)

			if err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to create dir for %s: %w", file.Path, err))
				mu.Unlock()
				return
			}

			// log.Println("Downloading:", file.Path)

			if err := downloadFile(file.DownloadUrl, localPath); err != nil {
				// log.Println("failed:", file.Path, err)
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to download %s: %w", file.Path, err))
				mu.Unlock()
				return
			}
		}()
	}

	// Wait for all downloads to finish
	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("some downloads failed: %v", errors)
	}

	// log.Println("Download complete")
	return nil
}

func handleDownloadWithoutDirTreePreserved(tree *tview.TreeView) error {
	var selected []*Content

	root := tree.GetRoot()
	if root == nil {
		return fmt.Errorf("tree is empty")
	}

	collectSelectedFromNode(root, &selected)

	if len(selected) == 0 {
		return fmt.Errorf("no files selected")
	}

	baseDir := reponame // flat download directory

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	for _, file := range selected {
		if file.IsDir {
			continue
		}

		wg.Add(1)

		// capture the variable for closure
		file := file

		go func() {
			defer wg.Done()

			// Save all files directly under baseDir, ignoring subdirectories
			localPath := filepath.Join(baseDir, filepath.Base(file.Path))

			if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to create dir for %s: %w", file.Path, err))
				mu.Unlock()
				return
			}

			if err := downloadFile(file.DownloadUrl, localPath); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to download %s: %w", file.Path, err))
				mu.Unlock()
				return
			}
		}()
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("some downloads failed: %v", errors)
	}

	return nil
}

func handleDownload_(tree *tview.TreeView) error {
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

		// log.Println("Downloading:", file.Path)

		if err := downloadFile(file.DownloadUrl, localPath); err != nil {
			// log.Println("failed:", file.Path, err)
			continue
		}
	}

	// log.Println("Download complete")
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

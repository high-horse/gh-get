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

	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error
	usedNames := make(map[string]int) // track used filenames

	for _, file := range selected {
		if file.IsDir {
			continue
		}

		wg.Add(1)
		file := file // capture for closure

		go func() {
			defer wg.Done()

			mu.Lock()
			name := filepath.Base(file.Path)
			count := usedNames[name]
			if count > 0 {
				ext := filepath.Ext(name)
				base := name[:len(name)-len(ext)]
				name = fmt.Sprintf("%s(%d)%s", base, count, ext)
			}
			usedNames[filepath.Base(file.Path)] = count + 1
			mu.Unlock()

			localPath := filepath.Join(baseDir, name)

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

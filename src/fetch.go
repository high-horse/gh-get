package main

import (
	"encoding/json"
	"fmt"
	// "log"
	
	"io"
	"net/http"
	"os"

)

func fetchContents(link string) error {
	mainLink, branchesLink, err := getUrls(link)
	if err != nil {
		return err
	}
	// log.Println("mainLink ", mainLink)
	// log.Println("branchesLink ", branchesLink)

	// Fetch repository info to get default_branch
	resp, err := hitHttpRequest(mainLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for non-success response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errMsg struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errMsg); err != nil {
			return fmt.Errorf("failed request with status %d", resp.StatusCode)
		}
		return fmt.Errorf("API error: %s", errMsg.Message)
	}

	var repo struct {
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return err
	}
	selectedBranch = repo.DefaultBranch

	// Fetch all branch names
	resp, err = hitHttpRequest(branchesLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errMsg struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errMsg); err != nil {
			return fmt.Errorf("failed request with status %d", resp.StatusCode)
		}
		return fmt.Errorf("API error: %s", errMsg.Message)
	}

	var branchList []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&branchList); err != nil {
		return err
	}

	branches = make([]string, len(branchList))
	for i, b := range branchList {
		branches[i] = b.Name
	}
	// log.Println("branches ", branches)
	// log.Println("selected branch 1  ", selectedBranch)

	return nil
}

func fetchContents_(link string) error {
	mainLink, branchesLink, err := getUrls(link)
	if err != nil {
		return err
	}
	// log.Println("mainLink ", mainLink)
	// log.Println("branchesLink ", branchesLink)

	// Fetch repository info to get default_branch
	resp, err := hitHttpRequest(mainLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var repo struct {
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return err
	}
	selectedBranch = repo.DefaultBranch

	// Fetch all branch names
	resp, err = hitHttpRequest(branchesLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var branchList []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&branchList); err != nil {
		return err
	}

	branches = make([]string, len(branchList))
	for i, b := range branchList {
		branches[i] = b.Name
	}
	// log.Println("branches ", branches)
	// log.Println("selected branch 1  ", selectedBranch)

	return nil
}

func fetchContentAtPath(owner, repo, branch, path string) ([]Content, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, path, branch)
	// log.Println("fetchContentAtPath url", url)

	resp, err := hitHttpRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw []struct {
		Name        string `json:"name"`
		Type        string `json:"type"` // "file" or "dir"
		DownloadUrl string `json:"download_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	contents := make([]Content, len(raw))
	for i, r := range raw {
		pathPrefix := path
		if pathPrefix != "" {
			pathPrefix += "/"
		}

		contents[i] = Content{
			Name:        r.Name,
			Path:        pathPrefix + r.Name, // ✅ FULL PATH
			IsDir:       r.Type == "dir",
			DownloadUrl: r.DownloadUrl,
			Children:    nil,
			Fetched:     false,
		}
		// contents[i] = Content{
		// 	Name:        r.Name,
		// 	IsDir:       r.Type == "dir",
		// 	DownloadUrl: r.DownloadUrl,
		// 	Children:    nil,
		// 	Fetched:     false,
		// }
	}
	return contents, nil
}

func fetchChildrenIfNeeded(c *Content, owner, repo, branch string) error {
	if c.Fetched || !c.IsDir {
		return nil
	}

	children, err := fetchContentAtPath(owner, repo, branch, c.Path)
	if err != nil {
		return err
	}
	c.Children = children
	c.Fetched = true
	return nil
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
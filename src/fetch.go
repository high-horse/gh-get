package main

import (
	"encoding/json"
	"log"
)

func fetchContents(link string) error {
	mainLink, branchesLink, err := getUrls(link)
	if err != nil {
		return err
	}

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
	log.Println("branches ", branches)
	log.Println("selected branch  ", selectedBranch)

	return nil
}

func fetchRepoContents(url string) error {
	return  nil
}
package main

type Content struct {
	Name        string
	IsDir       bool
	DownloadUrl string
}

var (
	url            string
	branches       []string
	selectedBranch string
	contents       []Content
	username       string
	reponame       string
)

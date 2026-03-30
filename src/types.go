package main

type Content struct {
	Name        string
	Path        string
	IsDir       bool
	DownloadUrl string
	Children    []Content
	Fetched     bool
	Selected    bool
}

var (
	url            string
	branches       []string
	selectedBranch string
	contents       []Content
	owner          string
	reponame       string
	preserveDirTree bool = true
)

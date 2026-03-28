package main

import (
	"fmt"
	"log"
	"strings"
)


func getUrls(link string) (string, string, error) {
	if !validateRepoLink(link) {
		return "", "", fmt.Errorf("invalid repo link URL: %s", link)
	}
	tokens, err := tokenize(link)
	if err != nil {
		log.Println("error ", err)
		return "" ,"", err
	}
	log.Println("checking contains " , link)
	if strings.Contains(link, "github.com") {
		// https://api.github.com/repos/torvalds/linux/branches
		username = tokens[3]
		reponame = tokens[4]
		mainLink := fmt.Sprintf("http://api.github.com/repos/%s/%s", tokens[3], tokens[4])
		branchLink := fmt.Sprintf("http://api.github.com/repos/%s/%s/branches", tokens[3], tokens[4])
		return mainLink, branchLink, nil
	} 
	return "","",  nil
}

func tokenize(link string) ([]string, error) {
	log.Println("Passed to tokenize ", link)
	tokens := strings.Split(link, "/")
	
	if len(tokens) < 3 {
		return nil, fmt.Errorf("invalid repo link URL: %s", link)
	}
	
	return tokens, nil
}
package main

import (
	// "log"
	"regexp"
	"strings"
	// "http/url"
	// "strings"
)

func validateRepoLink(link string) bool {
	strings.TrimSpace(link)
	pattern := `^https://github\.com/[\w\-]+/[\w\-]+/?$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(link)
	// return true
}

package main

import (
	// "log"
	"regexp"
	// "http/url"
	// "strings"
)

func validateRepoLink(link string) bool {
	pattern := `^https://github\.com/[\w\-]+/[\w\-]+/?$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(link)
	// return true
}

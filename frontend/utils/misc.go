package utils

import (
	"fmt"
	"regexp"
)

func ExtractFolderID(link string) (string, error) {
	// Define a regex pattern to match the Google Drive folder structure
	// and capture the folder ID
	pattern := `https://drive\.google\.com/drive/folders/([0-9A-Za-z-_]+)`
	re := regexp.MustCompile(pattern)

	// FindStringSubmatch returns a slice of strings holding the text of
	// the leftmost match of the regular expression in s and the matches,
	// if any, of its subexpressions, as defined by the 'SubexpNames' method.
	matches := re.FindStringSubmatch(link)

	if len(matches) > 1 {
		// The first match is the entire match, and the second is the captured group
		return matches[1], nil
	}

	return "", fmt.Errorf("folder ID not found in link")
}

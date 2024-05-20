package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/xhermitx/gitpulse-tracker/API"
	"github.com/xhermitx/gitpulse-tracker/models"
	"github.com/xhermitx/gitpulse-tracker/utils"
)

func getFileNames() ([]string, error) {
	matches, _ := filepath.Glob("*.pdf")
	if len(matches) == 0 {
		return nil, errors.New("no files in the current directory")
	}
	return matches, nil
}

func getUserName(fileName string) ([]string, error) {

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var contentBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		contentBuilder.WriteString(line)
	}

	content := contentBuilder.String()

	pattern := regexp.MustCompile(`https://github\.com/[a-zA-Z0-9]+(\-[a-zA-Z0-9]*)*`)

	uniqIDs := make(map[string]bool)

	// Find and print all matches
	matches := pattern.FindAllString(content, -1)
	for _, match := range matches {
		uniqIDs[match[19:]] = true
	}

	if len(uniqIDs) == 0 {
		return nil, fmt.Errorf("no username found in file : %s", f.Name())
	}

	userIDs := make([]string, 0, len(uniqIDs))

	for key := range uniqIDs {
		userIDs = append(userIDs, key)
	}

	return userIDs, nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading the environment variables")
	}

	fileNames, err := getFileNames()
	if err != nil {
		log.Println(err)
	}

	var userIDs []string

	for _, f := range fileNames {
		if ID, err := getUserName(f); err != nil {
			log.Println(err)
		} else {
			userIDs = append(userIDs, ID...)
		}
	}

	// fmt.Println("\nUSER IDs EXTRACTED: ", userIDs)

	var detailedList []models.GitResponse

	for _, user := range userIDs {
		if res, err := API.GetUserDetails(user); err != nil {
			log.Println(err)
		} else {
			detailedList = append(detailedList, res)
		}
	}

	utils.Printer(detailedList)

}

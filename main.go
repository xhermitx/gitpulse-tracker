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
	api "github.com/xhermitx/gitpulse-tracker/API"
)

func getFileNames() ([]string, error) {
	matches, _ := filepath.Glob("*.pdf")
	if len(matches) == 0 {
		return nil, errors.New("no files in the current directory")
	}
	return matches, nil
}

func getUserName(fileName string) ([]string, error) {
	// Replace 'resume.pdf' with the path to your PDF file
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Println("NAME of the file: ", f.Name())

	scanner := bufio.NewScanner(f)

	var contentBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		// Write the line to the header string builder
		contentBuilder.WriteString(line)
	}

	content := contentBuilder.String()

	pattern := regexp.MustCompile(`https://github\.com/[a-zA-Z0-9]+(\-[a-zA-Z0-9]*)*`)

	uniqIDs := make(map[string]bool)

	// Find and print all matches
	matches := pattern.FindAllString(content, -1)
	for _, match := range matches {
		// fmt.Println("GitHub Profile:", match[19:])
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

	var validUsers []string

	for _, name := range userIDs {
		if exists, error := api.CheckUserExists(name); error != nil {
			log.Println(err)
		} else if exists {
			validUsers = append(validUsers, name)
		}
	}

	fmt.Println(validUsers)

	userList := make(map[string]int)

	err = godotenv.Load()
	if err != nil {
		log.Panic("Error loading the environment variables")
	}

	for _, user := range validUsers {
		if contributions, err := api.GetContributions(user, os.Getenv("GITHUB_TOKEN")); err != nil {
			log.Print(err)
		} else {
			userList[user] = contributions
		}
	}

	for key, val := range userList {
		fmt.Printf("\nUSER : %s 	CONTRIBUTIONS : %d", key, val)
	}

}

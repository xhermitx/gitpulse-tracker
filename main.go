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
	"sync"
	"time"

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

	// PERFORMANCE CHECKS
	t := time.Now()
	defer func() {

		f, err := os.OpenFile("Performance.md", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Print(err)
		}
		defer f.Close()

		_, _ = fmt.Fprintln(f, "# PERFORMANCE WITH CONCURRENTLY READING PDFs AND FETCHING FROM GITHUB")
		_, _ = fmt.Fprintln(f, "TOTAL TIME TAKEN: ", time.Since(t).Seconds())
	}()

	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading the environment variables")
	}

	fileNames, err := getFileNames()
	if err != nil {
		log.Println(err)
	}

	var userIDs []string
	var lock sync.Mutex
	var wg sync.WaitGroup

	wg.Add(len(fileNames))
	for _, f := range fileNames {
		go func(f string) {
			defer wg.Done()
			if ID, err := getUserName(f); err != nil {
				log.Println(err)
			} else {
				lock.Lock()
				userIDs = append(userIDs, ID...)
				lock.Unlock()
			}
		}(f)
	}
	wg.Wait()

	fmt.Println("\nUSER IDs EXTRACTED: ", userIDs)

	var detailedList []models.GitResponse

	wg.Add(len(userIDs))
	for _, user := range userIDs {
		go func(user string) {

			defer wg.Done()

			if res, err := API.GetUserDetails(user); err != nil {
				log.Println(err)
			} else {
				lock.Lock()
				detailedList = append(detailedList, res)
				lock.Unlock()
			}
		}(user)
	}
	wg.Wait()

	utils.Printer(detailedList)

}

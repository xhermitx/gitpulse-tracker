package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	api "github.com/xhermitx/gitpulse-tracker/API"
)

func checkUserExists(username string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	// Perform the GET request.
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check the HTTP status code to determine if the user exists.
	if resp.StatusCode == http.StatusOK {
		var body interface{}
		decoder := json.NewDecoder(resp.Body)

		if err = decoder.Decode(&body); err != nil {
			log.Println(err)
			return false, err
		}

		// CHECK IF THE ENTITY IS OF TYPE "USER"
		return body.(map[string]interface{})["type"].(string) == "User", nil

	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("unexpected HTTP status: %s", resp.Status)

}

func main() {
	// Replace 'resume.pdf' with the path to your PDF file
	f, err := os.Open("resume.pdf")
	if err != nil {
		panic(err)
	}
	fmt.Println("NAME of the file: ", f.Name())

	scanner := bufio.NewScanner(f)

	var headerBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Write the line to the header string builder
		headerBuilder.WriteString(line)
	}

	content := headerBuilder.String()

	pattern := regexp.MustCompile(`https://github\.com/[a-zA-Z0-9]+(\-[a-zA-Z0-9]*)*`)

	userNames := make(map[string]bool)

	// Find and print all matches
	matches := pattern.FindAllString(content, -1)
	for _, match := range matches {
		// fmt.Println("GitHub Profile:", match[19:])
		userNames[match[19:]] = true
	}

	var validUsers []string

	for name := range userNames {
		if exists, error := checkUserExists(name); error != nil {
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

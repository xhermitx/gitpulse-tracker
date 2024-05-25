package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/xhermitx/gitpulse-tracker/API"
	"github.com/xhermitx/gitpulse-tracker/models"
	"github.com/xhermitx/gitpulse-tracker/utils"
)

func main() {

	// PERFORMANCE CHECKS
	t := time.Now()
	defer func() {

		f, err := os.OpenFile("Performance.md", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Print(err)
		}
		defer f.Close()
		_, _ = fmt.Fprintln(f, "# PERFORMANCE WITH CONCURRENTLY READING DRIVE DATA, PDFs AND FETCHING FROM GITHUB")
		_, _ = fmt.Fprintln(f, "TOTAL TIME TAKEN: ", time.Since(t).Seconds())
	}()

	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading the environment variables")
	}

	var lock sync.Mutex
	var wg sync.WaitGroup

	userIDs, err := API.GetDriveDetails()
	if err != nil {
		log.Fatal(err)
	}

	var detailedList []models.GitResponse

	// GET USER DETAILS FROM GITHUB
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

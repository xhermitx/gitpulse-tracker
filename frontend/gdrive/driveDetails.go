package gdrive

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func GetDriveDetails(folderID string) ([]string, error) {

	ctx := context.Background()

	// Read the service account key file
	log.Println(os.Getenv("CREDENTIALS_JSON"))

	data, err := os.ReadFile(os.Getenv("CREDENTIALS_JSON"))
	if err != nil {
		log.Println("Error reading Credentials")
		return nil, err
	}

	// CREATE CONFIGS USING AUTHENTICATION
	config, err := google.JWTConfigFromJSON(data, drive.DriveReadonlyScope)
	if err != nil {
		log.Println("Error getting JWT Configs")
		return nil, err
	}
	client := config.Client(ctx)

	// CREATE A NEW DRIVE CLIENT
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println("Error creating Service")
		return nil, err
	}

	// QUERY TO READ FILES RESIDING IN FOLDERS
	query := fmt.Sprintf("'%s' in parents", folderID)

	fileList, err := driveService.Files.List().
		Q(query).
		Fields("nextPageToken, files(id, name)").
		Do()
	if err != nil {
		log.Println("Error fetching the file list")
		return nil, err
	}

	var lock sync.Mutex
	var wg sync.WaitGroup
	var userIDs []string

	wg.Add(len(fileList.Files))
	// Print the names and IDs of the files
	for _, f := range fileList.Files {
		go func(f *drive.File) {
			defer wg.Done()
			// fmt.Printf("File Name: %s, File ID: %s\n", i.Name, i.Id)
			file, err := getFileDataWithRetry(f.Id, driveService)
			if err != nil {
				log.Print(err)
			}
			res, err := getUsername(file)
			if err != nil {
				log.Print(err)
			}

			log.Printf("\nFetched the following IDs from file %s : ", f.Name)
			log.Println(res)

			lock.Lock()
			userIDs = append(userIDs, res...)
			lock.Unlock()
		}(f)
	}

	wg.Wait()
	return userIDs, nil
}

func getFileDataWithRetry(fileID string, driveService *drive.Service) ([]byte, error) {
	var (
		content []byte
		err     error
	)

	for retries := 0; retries < 5; retries++ {
		content, err = getFileData(fileID, driveService)
		if err == nil {
			return content, nil
		}

		// Exponential backoff
		time.Sleep(time.Duration(retries) * time.Second)
	}
	return nil, err
}

func getFileData(fileID string, driveService *drive.Service) ([]byte, error) {

	// Get the file's metadata to verify it's a PDF
	file, err := driveService.Files.Get(fileID).Fields("mimeType").Do()
	if err != nil {
		return nil, err
	}

	// Only proceed if the file is a PDF
	if file.MimeType != "application/pdf" {
		return nil, err
	}

	// GET THE FILE'S CONTENT
	resp, err := driveService.Files.Get(fileID).Download()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the file content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func getUsername(data []byte) ([]string, error) {

	content := string(data)

	pattern := regexp.MustCompile(`github\.com/[a-zA-Z0-9]+(\-[a-zA-Z0-9]*)*`)

	uniqIDs := make(map[string]bool)

	// Find and print all matches
	matches := pattern.FindAllString(content, -1)
	for _, match := range matches {
		uniqIDs[match[11:]] = true
	}

	if len(uniqIDs) == 0 {
		return nil, fmt.Errorf("no username found in file")
	}

	userIDs := make([]string, 0, len(uniqIDs))

	for key := range uniqIDs {
		userIDs = append(userIDs, key)
	}

	return userIDs, nil
}

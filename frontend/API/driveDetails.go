package API

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func GetDriveDetails() ([]string, error) {
	ctx := context.Background()

	// Read the service account key file
	data, err := os.ReadFile(os.Getenv("CREDENTIALS_JSON"))
	if err != nil {
		return nil, err
	}

	// CREATE CONFIGS USING AUTHENTICATION
	config, err := google.JWTConfigFromJSON(data, drive.DriveReadonlyScope)
	if err != nil {
		return nil, err
	}
	client := config.Client(ctx)

	// CREATE A NEW DRIVE CLIENT
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	folderID := os.Getenv("FOLDER_ID")

	// QUERY TO READ FILES RESIDING IN FOLDERS
	query := fmt.Sprintf("'%s' in parents", folderID)

	fileList, err := driveService.Files.List().
		Q(query).
		Fields("nextPageToken, files(id, name)").
		Do()
	if err != nil {
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
			file, err := getFileData(f.Id, driveService)
			if err != nil {
				log.Print(err)
			}
			res, err := getUsername(file)
			if err != nil {
				log.Print(err)
			}
			lock.Lock()
			userIDs = append(userIDs, res...)
			lock.Unlock()
		}(f)
	}

	wg.Wait()
	return userIDs, nil
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

	pattern := regexp.MustCompile(`https://github\.com/[a-zA-Z0-9]+(\-[a-zA-Z0-9]*)*`)

	uniqIDs := make(map[string]bool)

	// Find and print all matches
	matches := pattern.FindAllString(content, -1)
	for _, match := range matches {
		uniqIDs[match[19:]] = true
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

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func CheckUserExists(username string) (bool, error) {
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

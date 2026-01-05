package net

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)


type LocalAuth struct {
	AuthToken string `json:"auth_token"`
}

var (
	Client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

func ValidTokenExists() bool {
	localData := GetSavedData()

	if localData.AuthToken == "" {
		return false
	}

	localAuth := LocalAuth{
		AuthToken: localData.AuthToken,
	}

	requestBody, err := json.Marshal(localAuth)
	if err != nil {
		fmt.Printf("An error occured while serializing request body: %v\n", err)
		return false
	}

	request, err := http.NewRequest("POST", "https://localhost:443/auth/valid", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("An error occured while constructing request: %v\n", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := Client.Do(request)
	if err != nil {
		fmt.Printf("An error occured while sending request: %v\n", err)
	}

	if response.StatusCode != 200 {
		err = RemoveLocalAuthToken()
		if err != nil {
			fmt.Println("Could now write to file")
		}

		return false
	}
	return true
}


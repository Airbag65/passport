package net

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)


func GetSavedData() *UserInformation {
	var data UserInformation

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	path := filepath.Join(homeDir, ".passport/authentication.json")

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil
	}

	defer func(){
		if err = jsonFile.Close(); err != nil {
			fmt.Println("Could not close file")
		}
	}()

	fileBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil
	}

	json.Unmarshal(fileBytes, &data)
	return &data
}

func RemoveLocalAuthToken() error {
	bytesToWrite, err := json.Marshal(UserInformation{})
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(homeDir, ".passport/authentication.json")
	os.WriteFile(path, bytesToWrite, 0644)
	return nil
}

func AddLocalAuthToken(authToken, name, surname, email string) error {
	bytesToWrite, err := json.Marshal(UserInformation{
		AuthToken: authToken,
		Name: name,
		Surname: surname,
		Email: email,
	})
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(homeDir, ".passport/authentication.json")
	os.WriteFile(path, bytesToWrite, 0644)
	return nil
}

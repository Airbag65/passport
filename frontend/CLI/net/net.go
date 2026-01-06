package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"passport-cli/enc"
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

func Login(email, password string) (*LoginResponse, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("No email or password provided\n")
	}
	loginRequestBody := LoginRequest{
		Email:    email,
		Password: password,
	}

	requestBodyBytes, err := json.Marshal(loginRequestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", "https://localhost:443/auth/login", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := Client.Do(request)
	if err != nil {
		return nil, err
	}

	var buffer []byte
	if response.StatusCode != 200 {
		return &LoginResponse{
			ResponseCode: response.StatusCode,
		}, nil
	}

	buffer, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var loginRes LoginResponse

	if err = json.Unmarshal(buffer, &loginRes); err != nil {
		return nil, err
	}

	if err = enc.StringToPEMFile(loginRes.PemString); err != nil {
		return nil, err
	}

	err = AddLocalAuthToken(loginRes.AuthToken, loginRes.Name, loginRes.Surname, loginRes.Email)
	if err != nil {
		return nil, err
	}

	return &loginRes, nil
}

func SignOut() (int, error) {
	email := GetSavedData().Email

	if email == "" {
		return 304, nil
	}
	signOutReq := SignOutRequest{
		Email: email,
	}

	reqBody, err := json.Marshal(signOutReq)
	if err != nil {
		return 0,err
	}

	request, err := http.NewRequest("PUT", "https://localhost:443/auth/signOut", bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, err
	}

	request.Header.Set("Content-Type", "application/json")

	res, err := Client.Do(request)
	if err != nil {
		return 0, err
	}

	if res.StatusCode == 200 {
		if err = RemoveLocalAuthToken(); err != nil {
			return 0, err
		}
	}

	return res.StatusCode, nil
}

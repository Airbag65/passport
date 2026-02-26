package main

import (
	"SH-password-manager/enc"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK!\n"))
}

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		BadRequest(w)
		return
	}

	var request LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		BadRequest(w)
		return
	}

	userInformation := s.GetUserWithEmail(request.Email)
	if userInformation == nil {
		NotFound(w)
		return
	}

	// if userInformation.AuthToken != "" {
	// 	w.WriteHeader(418)
	// 	w.Write([]byte("Already signed in"))
	// 	return
	// }

	if encryptPassword(request.Password) != userInformation.Password {
		Unauthorized(w)
		return
	}

	pemString, err := enc.PEMFileToString("publicKey")
	if err != nil {
		log.Printf("Error: %v", err)
		InternalServerError(w)
		return
	}

	loginResponse := &LoginResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		AuthToken:       "",
		Name:            userInformation.Name,
		Surname:         userInformation.Surname,
		Email:           userInformation.Email,
		PemString:       pemString,
	}
	newUserToken := fmt.Sprintf("%s%s", uuid.New().String(), uuid.New().String())
	newClientToken := fmt.Sprintf("%s%s", uuid.New().String(), uuid.New().String())
	ipAddr := GetRequestIP(r)

	s.SetNewAuthToken(request.Email, newUserToken, newClientToken, ipAddr)
	responseToken := newUserToken + "+" + newClientToken

	loginResponse.AuthToken = responseToken

	WriteJSON(w, loginResponse)
}

func (v *ValidateTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		BadRequest(w)
		return
	}

	var request ValidateTokenRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		BadRequest(w)
		return
	}

	userInformation := s.GetUserWithAuthToken(request.AuthToken)
	if userInformation == nil {
		Unauthorized(w)
		return
	}

	valid := s.ValidateToken(request.AuthToken, GetRequestIP(r), userInformation.Email)
	if !valid {
		Unauthorized(w)
		return
	}

	pemString, err := enc.PEMFileToString("publicKey")
	if err != nil {
		InternalServerError(w)
		return
	}

	tokenResponse := &ValidateTokenResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		Name:            userInformation.Name,
		Surname:         userInformation.Surname,
		Email:           userInformation.Email,
		PemString:       pemString,
	}
	WriteJSON(w, tokenResponse)
}

func (h *SignOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		MethodNotAllowed(w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		BadRequest(w)
		return
	}

	var request SignOutRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		BadRequest(w)
		return
	}

	user := s.GetUserWithEmail(request.Email)
	if user == nil {
		NotFound(w)
		return
	}

	if user.AuthToken == "" {
		w.WriteHeader(304)
		w.Write([]byte("Not modified"))
		return
	}

	s.RemoveAuthToken(request.Email, GetRequestIP(r))

	WriteJSON(w, SignOutResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
	})
}

func (c *CreateNewUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		BadRequest(w)
		return
	}

	var request CreateNewUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		BadRequest(w)
		return
	}

	existingUser := s.GetUserWithEmail(request.Email)
	if existingUser != nil {
		w.WriteHeader(418)
		w.Write([]byte("User already exists"))
		return
	}

	pemString, err := enc.PEMFileToString("publicKey")
	if err != nil {
		InternalServerError(w)
		return
	}

	encPwd := encryptPassword(request.Password)

	s.CreateNewUser(request.Email, encPwd, request.Name, request.Surname)

	newUserResponse := &CreateNewUserResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		AuthToken:       "",
		Name:            request.Name,
		Surname:         request.Surname,
		PemString:       pemString,
	}
	newUserToken := fmt.Sprintf("%s%s", uuid.New().String(), uuid.New().String())
	newClientToken := fmt.Sprintf("%s%s", uuid.New().String(), uuid.New().String())
	ipAddr := GetRequestIP(r)
	s.SetNewAuthToken(request.Email, newUserToken, newClientToken, ipAddr)
	responseToken := newUserToken + "+" + newClientToken
	newUserResponse.AuthToken = responseToken
	WriteJSON(w, newUserResponse)
}

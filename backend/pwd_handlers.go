package main

import (
	"SH-password-manager/enc"
	"encoding/json"
	"log"
	"net/http"
)

func (h *GetPasswordHostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		MethodNotAllowed(w)
		return
	}

	userInformation := ValidateToken(w, r)

	names := s.GetHostNames(userInformation.Id)

	WriteJSON(w, GetPasswordHostsResponse{Hosts: names})
}

func (h *UploadNewPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w)
		return
	}

	var request UploadNewPasswordRequest

	if r.Header.Get("Content-Type") != "application/json" {
		BadRequest(w)
		return
	}
	userInformation := ValidateToken(w, r)

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		BadRequest(w)
		return
	}

	err = s.AddNewPassord(userInformation.Id, request.Password, request.HostName)
	if err != nil {
		InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (h *GetPasswordValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		MethodNotAllowed(w)
		return
	}

	userInformation := ValidateToken(w, r)

	var request GetPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		BadRequest(w)
		return
	}

	encPassword, err := s.GetPassword(userInformation.Id, request.HostName)
	if err != nil {
		NotFound(w)
		return
	}

	privatePemString, err := enc.PEMFileToString("privateKey")
	if err != nil {
		InternalServerError(w)
		return
	}

	privateKey := enc.PemStringToPrivateKey(privatePemString)

	decPassword, err := enc.Decrypt(encPassword, privateKey)
	if err != nil {
		log.Printf("Error: %v", err)
		InternalServerError(w)
		return
	}

	WriteJSON(w, GetPasswordResonse{Password: decPassword})
}

func (h *RemovePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		MethodNotAllowed(w)
		return
	}

	userInformation := ValidateToken(w, r)
	var req RemovePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		InternalServerError(w)
	}

	if err := s.RemovePassword(userInformation.Id, req.HostName); err != nil {
		NotFound(w)
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (h *EditPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		MethodNotAllowed(w)
		return
	}

	userInformation := ValidateToken(w, r)
	var req EditPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		InternalServerError(w)
	}

	if err := s.EditPassword(userInformation.Id, req.HostName, req.NewPassword); err != nil {
		NotFound(w)
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

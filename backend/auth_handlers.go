package main

import (
	"SH-password-manager/enc"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK!\n"))
}

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	version := GetRustportVersion()
	if version == "" {
		InternalServerError(w)
		return
	}
	WriteJSON(w, &StatusResponse{Health: "OK!", RustportVersion: version})
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

	authportReqObj := &AuthportLoginRequest{
		Email:            request.Email,
		Password:         encryptPassword(request.Password),
		ClientIdentifier: "cli",
		RemoteAddr:       GetRequestIP(r),
	}

	authportReqBody, err := json.Marshal(authportReqObj)
	if err != nil {
		InternalServerError(w)
		return
	}

	authportReq, err := http.NewRequest("POST", "http://127.0.0.1:8000/login", bytes.NewBuffer(authportReqBody))
	if err != nil {
		InternalServerError(w)
		return
	}
	authportReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(authportReq)
	if err != nil {
		InternalServerError(w)
		return
	}
	switch response.StatusCode {
	case 401:
		Unauthorized(w)
		return
	case 404:
		NotFound(w)
		return
	case 418:
		w.WriteHeader(418)
		w.Write([]byte("Already logged in"))
		return
	}
	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		InternalServerError(w)
		return
	}

	var authportResponse AuthportLoginResponse

	if err = json.Unmarshal(buffer, &authportResponse); err != nil {
		InternalServerError(w)
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
		AuthToken:       authportResponse.AuthToken,
		Name:            authportResponse.Name,
		Surname:         authportResponse.Surname,
		Email:           request.Email,
		PemString:       pemString,
	}

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

	if request.AuthToken == "" {
		Unauthorized(w)
		return
	}

	authportReqBody, err := json.Marshal(request)
	if err != nil {
		InternalServerError(w)
		return
	}

	authportReq, err := http.NewRequest("POST", "http://127.0.0.1:8000/valid", bytes.NewBuffer(authportReqBody))
	if err != nil {
		InternalServerError(w)
		return
	}

	authportReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(authportReq)
	if err != nil {
		InternalServerError(w)
		return
	}

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		InternalServerError(w)
		return
	}

	var authportResponse AuthportValidateTokenResponse

	if err = json.Unmarshal(buffer, &authportResponse); err != nil {
		InternalServerError(w)
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
		Name:            authportResponse.Name,
		Surname:         authportResponse.Surname,
		Email:           authportResponse.Email,
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

	authportReqObj := &AuthportSignoutRequest{
		Email:            request.Email,
		IpAddr:           GetRequestIP(r),
		ClientIdentifier: "cli",
	}
	authportReqBody, err := json.Marshal(authportReqObj)
	if err != nil {
		InternalServerError(w)
		return
	}

	authportReq, err := http.NewRequest("PUT", "http://127.0.0.1:8000/signOut", bytes.NewBuffer(authportReqBody))
	if err != nil {
		InternalServerError(w)
		return
	}
	authportReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(authportReq)

	switch response.StatusCode {
	case 304:
		w.WriteHeader(304)
		w.Write([]byte("Not modified"))
		return
	case 404:
		NotFound(w)
		return
	}

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

	request.Password = encryptPassword(request.Password)

	authportReqBody, err := json.Marshal(request)
	if err != nil {
		InternalServerError(w)
		return
	}

	client := &http.Client{}
	authportReq, err := http.NewRequest("POST", "http://127.0.0.1:8000/new", bytes.NewBuffer(authportReqBody))
	if err != nil {
		InternalServerError(w)
		return
	}

	authportReq.Header.Set("Content-Type", "application/json")
	response, err := client.Do(authportReq)
	if response.StatusCode == 418 {
		w.WriteHeader(418)
		w.Write([]byte("User already exists"))
		return
	} else if response.StatusCode != 200 {
		InternalServerError(w)
		return
	}

	authportLoginReqObj := &AuthportLoginRequest{
		Email:            request.Email,
		Password:         request.Password,
		ClientIdentifier: "cli",
		RemoteAddr:       GetRequestIP(r),
	}

	authportLoginReqBody, err := json.Marshal(authportLoginReqObj)
	if err != nil {
		InternalServerError(w)
		return
	}

	authportReq, err = http.NewRequest("POST", "http://127.0.0.1:8000/login", bytes.NewBuffer(authportLoginReqBody))
	authportReq.Header.Set("Content-Type", "application/json")
	response, err = client.Do(authportReq)
	if err != nil {
		InternalServerError(w)
		return
	}

	switch response.StatusCode {
	case 404:
		NotFound(w)
		return
	case 401:
		Unauthorized(w)
		return
	case 418:
		w.WriteHeader(418)
		w.Write([]byte("Already logged in"))
		return
	case 500:
		InternalServerError(w)
		return
	}

	pemString, err := enc.PEMFileToString("publicKey")
	if err != nil {
		InternalServerError(w)
		return
	}

	var loginRes AuthportLoginResponse

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		InternalServerError(w)
		return
	}
	if err = json.Unmarshal(buffer, &loginRes); err != nil {
		InternalServerError(w)
		return
	}

	newUserResponse := &CreateNewUserResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		AuthToken:       loginRes.AuthToken,
		Name:            loginRes.Name,
		Surname:         loginRes.Surname,
		PemString:       pemString,
	}
	WriteJSON(w, newUserResponse)
}

func (h *RequestResetAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		MethodNotAllowed(w)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		BadRequest(w)
		return
	}

	var request RequestResetAccountRequest
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
	tokenString, err := CreateToken(request.Email, user.Name, user.Surname)
	if err != nil {
		InternalServerError(w)
		return
	}
	var usedIP string
	if request.Debug {
		usedIP = "127.0.0.1"
	} else {
		usedIP = GetLocalIP()
	}
	url := fmt.Sprintf("https://%s:443/auth/resetAcc?token=%s", usedIP, tokenString)
	err = SendResetEmail(url, request.Email, user.Name, user.Surname)
	if err != nil {
		InternalServerError(w)
		return
	}
	res := &RequsetResetAccountResponse{
		Url: url,
	}
	WriteJSON(w, res)
}

func (h *ResetAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		MethodNotAllowed(w)
	}
	ctx := r.Context()
	val := ctx.Value("token").(string)
	token, err := VerifyToken(val)
	if err != nil {
		InternalServerError(w)
		return
	}
	email, _, _ := ParseToken(token)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, ReadHTMLFile("reset.html"), email, val)
}

func (h *AccountResetter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w)
		return
	}
	formValues := map[string]string{
		"new_pass": r.FormValue("new_password"),
		"token":    r.FormValue("token"),
	}

	token, err := VerifyToken(formValues["token"])
	if err != nil {
		http.Redirect(w, r, "/auth/account/reset/error", http.StatusPermanentRedirect)
		return
	}
	email, _, _ := ParseToken(token)
	if email == "" {
		http.Redirect(w, r, "/auth/account/reset/error", http.StatusPermanentRedirect)
		return
	}

	user := s.GetUserWithEmail(email)
	if user == nil {
		http.Redirect(w, r, "/auth/account/reset/error", http.StatusPermanentRedirect)
		return
	}
	err = s.UpdateAccountPassword(user.Email, encryptPassword(formValues["new_pass"]))
	if err != nil {
		http.Redirect(w, r, "/auth/account/reset/error", http.StatusPermanentRedirect)
		return
	}
	err = s.SignOutUserAll(user.Email)
	if err != nil {
		http.Redirect(w, r, "/auth/account/reset/error", http.StatusPermanentRedirect)
		return
	}
	http.Redirect(w, r, "/auth/account/reset/success", http.StatusPermanentRedirect)
}

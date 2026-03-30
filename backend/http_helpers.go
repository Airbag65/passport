package main

import (
	"SH-password-manager/db"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Write([]byte("Bad Request"))
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("Internal Server Error"))
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte("Not Found"))
}

func MethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(405)
	w.Write([]byte("Method Not Allowed"))
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(401)
	w.Write([]byte("Unauthorized"))
}

func WriteJSON(w http.ResponseWriter, v any) error {
	// w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(v)
}

func ValidateToken(w http.ResponseWriter, r *http.Request) *db.User {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		BadRequest(w)
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/validate", &bytes.Buffer{})
	if err != nil {
		InternalServerError(w)
		return nil
	}
	req.Header.Set("Authorization", tokenHeader)
	response, err := client.Do(req)
	if err != nil {
		InternalServerError(w)
		return nil
	}

	switch response.StatusCode {
	case 401:
		Unauthorized(w)
		return nil
	case 400:
		BadRequest(w)
		return nil
	}

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		InternalServerError(w)
		return nil
	}
	var userInformation db.User

	if err = json.Unmarshal(buffer, &userInformation); err != nil {
		InternalServerError(w)
		return nil
	}

	if &userInformation == nil {
		Unauthorized(w)
		return nil
	}
	return &userInformation
}

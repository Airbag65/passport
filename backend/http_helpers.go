package main

import (
	"SH-password-manager/db"
	"encoding/json"
	"net/http"
	"strings"
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

	bearer := strings.Split(tokenHeader, " ")[0]
	if bearer != "Bearer" {
		BadRequest(w)
		return nil
	}
	token := strings.Split(tokenHeader, " ")[1]

	userInformation := s.GetUserWithAuthToken(token)
	if userInformation == nil {
		Unauthorized(w)
		return nil
	}
	return userInformation
}

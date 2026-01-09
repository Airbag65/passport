package net

import (
	"crypto/tls"
	"net/http"
)

var (
	Client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

type UserInformation struct {
	AuthToken string `json:"auth_token"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
}

type SignOutRequest struct {
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	AuthToken       string `json:"auth_token"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Email           string `json:"email"`
	PemString       string `json:"pem_string"`
}

type LocalAuth struct {
	AuthToken string `json:"auth_token"`
}

type Hosts struct {
	Hosts []string `json:"hosts"`
}

type SignupResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	AuthToken       string `json:"auth_token"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	PemString       string `json:"pem_string"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type getPasswordRequest struct {
	HostName string `json:"host_name"`
}

type getPasswordResonse struct {
	Password string `json:"password"`
}

type createPasswordRequest struct {
	HostName string `json:"host_name"`
	Password string `json:"password"`
}

type DeletePasswordRequest struct {
	HostName string `json:"host_name"`
}

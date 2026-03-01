package main

type HomeHandler struct{}

/*
--- LOGIN ---
*/
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

type LoginHandler struct{}

/*
--- VALIDATE TOKEN ---
*/
type ValidateTokenHandler struct{}

type ValidateTokenRequest struct {
	AuthToken string `json:"auth_token"`
}

type ValidateTokenResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Email           string `json:"email"`
	PemString       string `json:"pem_string"`
}

/*
   --- SIGN OUT ---
*/

type SignOutHandler struct{}

type SignOutRequest struct {
	Email string `json:"email"`
}

type SignOutResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
}

/*
--- CREATE NEW USER ---
*/
type CreateNewUserHandler struct{}

type CreateNewUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type CreateNewUserResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	AuthToken       string `json:"auth_token"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	PemString       string `json:"pem_string"`
}

type UploadNewPasswordHandler struct{}

type UploadNewPasswordRequest struct {
	HostName string `json:"host_name"`
	Password string `json:"password"`
}

type GetPasswordValueHandler struct{}

type GetPasswordRequest struct {
	HostName string `json:"host_name"`
}

type GetPasswordResonse struct {
	Password string `json:"password"`
}

type GetPasswordHostsHandler struct{}

type GetPasswordHostsResponse struct {
	Hosts []string `json:"hosts"`
}

type RemovePasswordHandler struct{}

type RemovePasswordRequest struct {
	HostName string `json:"host_name"`
}

type EditPasswordHandler struct{}

type EditPasswordRequest struct {
	HostName    string `json:"host_name"`
	NewPassword string `json:"new_password"`
}

type RequestResetAccountHandler struct{}

type RequestResetAccountRequest struct {
	Email string `json:"email"`
	Debug bool   `json:"debug,omitempty"`
}

type RequsetResetAccountResponse struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

type ResetAccountHandler struct{}

type ctxKey struct{}

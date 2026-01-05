package net


type UserInformation struct {
	AuthToken string `json:"auth_token"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
}


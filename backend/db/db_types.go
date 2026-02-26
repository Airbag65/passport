package db

import (
	"database/sql"
	"fmt"
)

type User struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Id            int    `json:"id"`
	AuthToken     string `json:"auth_token"`
	LoggedInCount int    `json:"logged_in_count"`
}

type AuthKey struct {
	Id              int    `json:"id"`
	IpAddr          string `json:"ip_addr"`
	UserToken       string `json:"user_token"`
	ClientToken     string `json:"client_token"`
	TokenExpiryDate int64  `json:"token_expiry_date"`
}

func (user *User) ToString() string {
	return fmt.Sprintf("Email: %s\nPassword: %s\nName: %s\nSurname: %s\nId: %d\nAuthToken: %s\n",
		user.Email,
		user.Password,
		user.Name,
		user.Surname,
		user.Id,
		user.AuthToken)
}

func DbEntryToUser(row *sql.Rows) *User {
	selectedUser := &User{}
	for row.Next() {
		row.Scan(&selectedUser.Id, &selectedUser.Email, &selectedUser.Password, &selectedUser.Name, &selectedUser.Surname, &selectedUser.AuthToken, &selectedUser.LoggedInCount)
	}
	if selectedUser.Name == "" {
		return nil
	}
	return selectedUser
}

func DbEntryToAuthKey(row *sql.Rows) *AuthKey {
	authKey := &AuthKey{}
	for row.Next() {
		row.Scan(&authKey.Id, &authKey.IpAddr, &authKey.UserToken, &authKey.ClientToken, &authKey.TokenExpiryDate)
	}
	if authKey.UserToken == "" {
		return nil
	}
	return authKey
}
func DbEntryToHostNames(rows *sql.Rows) []string {
	hostNames := []string{}
	for rows.Next() {
		name := new(string)
		rows.Scan(&name)
		hostNames = append(hostNames, *name)
	}
	return hostNames
}

func DbEntryToPassword(row *sql.Rows) string {
	var password string
	for row.Next() {
		row.Scan(&password)
	}

	return password
}

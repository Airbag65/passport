package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func extractSalt(pwd string) (string, string) {
	saltBeginning := ""
	for i := len(pwd) - 2; i < len(pwd); i++ {
		saltBeginning = fmt.Sprintf("%s%c", saltBeginning, pwd[i])
	}

	saltEnd := ""
	for _, c := range pwd[:5] {
		saltEnd = fmt.Sprintf("%s%c", saltEnd, c)
	}
	return saltBeginning, saltEnd
}

func encryptPassword(origPwd string) string {
	encPassword := origPwd
	saltBeginning, saltEnd := extractSalt(encPassword)

	sha256Encoder := sha256.New()
	sha256Encoder.Write([]byte(encPassword))
	encPassword = fmt.Sprintf("%x", sha256Encoder.Sum(nil))

	encPassword = fmt.Sprintf("%s%s%s", saltBeginning, encPassword, saltEnd)

	sha256Encoder.Write([]byte(encPassword))
	encPassword = fmt.Sprintf("%x", sha256Encoder.Sum(nil))
	return encPassword
}

func GetJWTSecret() []byte {
	if err := godotenv.Load(".env"); err != nil {
		return []byte{}
	}
	secretString := os.Getenv("JWT_SECRET")
	return []byte(secretString)
}

func CreateToken(email, name, surname string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"iss": "passport",
		"aud": fmt.Sprintf("%s %s", name, surname),
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	secret := GetJWTSecret()
	if string(secret) == "" {
		return "", fmt.Errorf("Could not retrieve JTW secret")
	}
	tokenString, err := claims.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid Token")
	}

	return token, nil
}

func ParseToken(token *jwt.Token) (string, string, string) {
	aud, err := token.Claims.GetAudience()
	if err != nil {
		return "", "", ""
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return "", "", ""
	}
	audience := strings.Split(strings.Join(aud, " "), " ")
	return sub, audience[0], audience[1]
}

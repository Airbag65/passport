package main

import (
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func GetRequestIP(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func GetSMTPPassword() string {
	if err := godotenv.Load(".env"); err != nil {
		return ""
	}
	smtpPass := os.Getenv("SMTP_PW")
	return smtpPass
}

func GetSMTPUsername() string {
	if err := godotenv.Load(".env"); err != nil {
		return ""
	}
	smtpPass := os.Getenv("SMTP_UN")
	return smtpPass
}

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:88")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func GetQueryParams(url *url.URL) map[string]string {
	fullParamString := strings.Split(url.String(), "?")[1]
	params := strings.Split(fullParamString, "&")
	result := make(map[string]string)
	for _, param := range params {
		item := strings.Split(param, "=")
		result[item[0]] = item[1]
	}
	return result
}

func SendResetEmail(url, email, name, surname string) error {
	auth := smtp.PlainAuth("", GetSMTPUsername(), GetSMTPPassword(), "smtp.gmail.com")

	to := []string{email}
	msg := []byte("From: \"PASSPORT\" <passport@noreply.com>\r\n" +
		fmt.Sprintf("To: \"%s %s\" <%s>\r\n", name, surname, email) +
		"Subject: [PASSPORT] Account reset was requested\r\n" +
		fmt.Sprintf("Hello %s %s,\r\nYou have requested to reset the master password for your PASSPORT account. You can do so via this link: %s \r\n\r\nPASSPORT wishes you a great day!",
			name, surname, url))

	err := smtp.SendMail("smtp.gmail.com:587", auth, "passport", to, msg)
	if err != nil {
		return err
	}
	return nil
}

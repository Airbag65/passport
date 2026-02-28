package main

import (
	"SH-password-manager/db"
	"SH-password-manager/enc"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

var s db.Storage

func main() {
	s = db.NewLocalStorage()
	if err := s.Init(); err != nil {
		panic("Could not initialize database")
	}
	f, err := os.OpenFile("target/log/file.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Println("Could not close log-file")
		}
	}()

	log.SetOutput(f)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			s.Migrate()
			return
		case "keygen":
			fmt.Println("Generating Keys")
			private, public, err := enc.GenerateKeys()
			if err != nil {
				return
			}
			fmt.Println("Making PEM strings")
			privatePemString := enc.PrivateKeyToPemString(private)
			publicPemString := enc.PublicKeyToPemString(public)
			fmt.Println("Saving PEM files")
			if err = enc.StringToPEMFile(privatePemString, "privateKey"); err != nil {
				return
			}
			if err = enc.StringToPEMFile(publicPemString, "publicKey"); err != nil {
				return
			}
		default:
			return
		}
	}

	server := http.NewServeMux()

	// Auth handlers
	server.Handle("/", &HomeHandler{})
	server.Handle("/auth/login", &LoginHandler{})
	server.Handle("/auth/valid", &ValidateTokenHandler{})
	server.Handle("/auth/signOut", &SignOutHandler{})
	server.Handle("/auth/new", &CreateNewUserHandler{})
	server.Handle("/auth/reset", &RequestResetAccountHandler{})
	server.Handle("/auth/reset/:jwt", &ResetAccountHandler{})

	// PWD handlers
	server.Handle("/pwd/getHosts", &GetPasswordHostsHandler{})
	server.Handle("/pwd/new", &UploadNewPasswordHandler{})
	server.Handle("/pwd/get", &GetPasswordValueHandler{})
	server.Handle("/pwd/remove", &RemovePasswordHandler{})
	server.Handle("/pwd/edit", &EditPasswordHandler{})

	handler := cors.Default().Handler(server)

	fmt.Println("Server running on: https://localhost:443...")
	err = http.ListenAndServeTLS(":443", "cert.pem", "key.pem", handler)
	if err != nil {
		log.Println("Could not start server")
		log.Fatal(err)
	}
}

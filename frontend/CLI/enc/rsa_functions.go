package enc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
)


func PublicKeyToPemString(key *rsa.PublicKey) string {
	keyDER := x509.MarshalPKCS1PublicKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "PUBLIC KEY",
		Bytes: keyDER,
	})
	return string(keyPEM)
}

func PemStringToPublicKey(pemString string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(pemString))
	parseRes, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return parseRes
}

func StringToPEMFile(pemString string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(homeDir, ".passport/publicKey.pem")

	os.WriteFile(path, []byte(pemString), 0644)
	return nil
}

func PEMFileToString() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homeDir, ".passport/publicKey.pem")

	pemFile, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer func(){
		if err = pemFile.Close(); err != nil {
			fmt.Println("Could not close file")
		}
	}()

	pemBytes, err := io.ReadAll(pemFile)
	if err != nil {
		return "", err
	}
	
	return string(pemBytes), nil
}

func Encrypt(message string, publicKey *rsa.PublicKey) (string, error) {
	encBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(message),
		nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encBytes), nil
}

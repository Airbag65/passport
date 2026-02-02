package enc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
)

func GenerateKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publicKey := privateKey.PublicKey
	return privateKey, &publicKey, nil
}

func PrivateKeyToPemString(key *rsa.PrivateKey) string {
	keyDER := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyDER,
	})
	return string(keyPEM)
}

func PemStringToPrivateKey(pemString string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemString))
	parseRes, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return parseRes
}

func PublicKeyToPemString(key *rsa.PublicKey) string {
	keyDER := x509.MarshalPKCS1PublicKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyDER,
	})
	return string(keyPEM)
}

func StringToPEMFile(pemString, fileName string) error {
	path := fmt.Sprintf("./target/pem/%s.pem", fileName)

	os.WriteFile(path, []byte(pemString), 0644)
	return nil
}

func PEMFileToString(fileName string) (string, error) {

	path := fmt.Sprintf("./target/pem/%s.pem", fileName)

	pemFile, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer func() {
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

func PemStringToPublicKey(pemString string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(pemString))
	parseRes, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return parseRes
}

func Encrypt(message string, publicKey *rsa.PublicKey) ([]byte, error) {
	encBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(message),
		nil)
	if err != nil {
		log.Println("Could not Encrypt")
		return []byte(""), err
	}
	return encBytes, nil
}

func Decrypt(encryptedMessage string, privateKey *rsa.PrivateKey) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
	decMessage, err := privateKey.Decrypt(nil, bytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		log.Println("Could not Decrypt")
		return "", err
	}
	return string(decMessage), nil
}

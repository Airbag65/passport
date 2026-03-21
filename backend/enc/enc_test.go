package enc_test

import (
	"SH-password-manager/enc"
	"fmt"
	"testing"
)

func TestKeys(t *testing.T) {
	private, public, err := enc.GenerateKeys()
	if err != nil {
		t.Error("Could not generate keys")
	}

	privateString := enc.PrivateKeyToPemString(private)
	newPrivate := enc.PemStringToPrivateKey(privateString)
	if fmt.Sprintf("%d", private.D) != fmt.Sprintf("%d", newPrivate.D) {
		t.Error("Private key changed after PEM")
	}

	publicString := enc.PublicKeyToPemString(public)
	newPublic := enc.PemStringToPublicKey(publicString)
	if public.E != newPublic.E {
		t.Error("Public key changed after PEM")
	}
}

func TestEncrypt(t *testing.T) {
	_, public, err := enc.GenerateKeys()
	if err != nil {
		t.Error("Could not generate keys")
	}

	_, err = enc.Encrypt("Super secret message", public)
	if err != nil {
		t.Error("Error while encrypting")
	}
}

func TestDecrypt(t *testing.T) {
	private, public, err := enc.GenerateKeys()
	if err != nil {
		t.Error("Could not generate keys")
	}
	message := "super secret message"

	encMsg, _ := enc.Encrypt(message, public)
	decMsg, err := enc.Decrypt(string(encMsg), private)
	if err != nil {
		t.Error("Error while decrypting")
	}
	if message != decMsg {
		t.Error("Decrypted message does not match original")
	}
}

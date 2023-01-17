package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io/ioutil"
	"testing"
)

var privateKey *rsa.PrivateKey

func TestGenerateKeys(t *testing.T) {
	_, err := GenerateKeys()
	if err != nil {
		t.Errorf("Error while generating keys: %v\n", err)
	}
}

func TestEncryptOAEP(t *testing.T) {
	keys, err := GenerateKeys()
	if err != nil {
		t.Errorf("Error while generating keys: %v\n", err)
	}

	privateKey = keys

	encrypt, err := EncryptOAEP(sha256.New(), rand.Reader, &keys.PublicKey, []byte("hello world"))
	if err != nil {
		t.Errorf("Error while calling EncryptOAEP(): %v\n", err)
	}

	if err := ioutil.WriteFile("encrypted.txt", encrypt, 0644); err != nil {
		t.Errorf("Error while writing encrypted data to file: %v\n", err)
	}
}

func TestDecryption(t *testing.T) {
	data, err := ioutil.ReadFile("encrypted.txt")
	if err != nil {
		t.Errorf("Error while reading encrypted data from file: %v\n", err)
	}

	decryptedData, err := DecryptOAEP(sha256.New(), rand.Reader, privateKey, data)
	if err != nil {
		t.Errorf("Error while calling DecryptOAEP(): %v\n", err)
	}

	if string(decryptedData) != "hello world" {
		t.Errorf("Unexpected result. Got %v want %v\n", string(decryptedData), "hello world")
	}
}

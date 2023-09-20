package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

/*
	Galois/Counter Mode (GCM) is a mode of operation for symmetric-key cryptographic
*/
func EncryptAESGCM(key string, plainText string) (encryptedText string, err error) {

	plaintext := []byte(plainText)

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	// NewGCM = GCM 128 bit
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	// Base 16
	// encryptedText = fmt.Sprintf("%x", ciphertext)

	// Base 64
	encryptedText = base64.StdEncoding.EncodeToString(ciphertext)

	return
}

func DecryptAESGCM(key string, encryptedText string) (plainText string, err error) {

	// Base 16
	// encryptedData, _ := hex.DecodeString(encryptedText)

	// Base 64
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedText)

	if err != nil {
		return
	}

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt the data
	plainTextBytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return
	}

	plainText = fmt.Sprintf("%s", plainTextBytes)

	return
}

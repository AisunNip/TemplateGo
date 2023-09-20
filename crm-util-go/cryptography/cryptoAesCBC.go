package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func EncryptAESCBC(key string, plainText string, iv []byte) (encryptedText string, err error) {
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return encryptedText, err
	}

	if plainText == "" {
		return encryptedText, err
	}

	cbcBlockMode := cipher.NewCBCEncrypter(block, iv)
	binaryPlainText := []byte(plainText)
	binaryPlainText = PKCS5Padding(binaryPlainText, block.BlockSize())

	encrypted := make([]byte, len(binaryPlainText))
	cbcBlockMode.CryptBlocks(encrypted, binaryPlainText)

	encryptedText = base64.StdEncoding.EncodeToString(encrypted)

	return encryptedText, nil
}

func DecryptAESCBC(key string, encryptedText string, iv []byte) (plainText string, err error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedText)

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return plainText, err
	}

	cbcBlockMode := cipher.NewCBCDecrypter(block, iv)
	binaryPlainText := make([]byte, len(encrypted))
	cbcBlockMode.CryptBlocks(binaryPlainText, encrypted)
	binaryPlainText = PKCS5Unpadding(binaryPlainText)

	plainText = string(binaryPlainText)

	return plainText, err
}

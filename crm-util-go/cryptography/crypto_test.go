package cryptography_test

import (
	"crm-util-go/cryptography"
	"fmt"
	"testing"
)

func TestAESGcm(t *testing.T) {

	key := "@TRUE-IT^!#$+*&%"
	data := "<Paravit#123/ทำการทดสอบ!!>"

	encryptedText, err := cryptography.EncryptAESGCM(key, data)

	if err != nil {
		t.Errorf("TestAESGcm - EncryptAESGCM Error %v", err.Error())
	} else {
		fmt.Println("EncryptedText:", encryptedText)

		decryptedText, err := cryptography.DecryptAESGCM(key, encryptedText)

		if err != nil {
			t.Errorf("TestAESGcm - DecryptAESGCM Error %v", err.Error())
		} else {
			fmt.Println("DecryptedText:", decryptedText)
		}
	}
}

func TestAES128Ecb(t *testing.T) {
	/*
	 * src string to be encrypted
	 * key Key used for encryption Key length can be any one of 128bit, 192bit, 256bit
	 * 16-bit key corresponds to 128bit
	 */
	key := "@TRUE-IT^!#$+*&%"
	data := "<Paravit#123/ทำการทดสอบ!!>"

	encryptedText, err := cryptography.EncryptAESECB(key, data)

	if err != nil {
		t.Errorf("TestAES128Ecb - EncryptAESECB Error %v", err.Error())
	} else {
		fmt.Println("encryptedText:", encryptedText)
	}

	plainText, err := cryptography.DecryptAESECB(key, encryptedText)

	if err != nil {
		t.Errorf("TestAES128Ecb - DecryptAESECB Error %v", err.Error())
	} else if data != plainText {
		t.Errorf("data not equal plainText")
	} else {
		fmt.Println("plainText:", plainText)
	}
}

func TestAES256Ecb(t *testing.T) {
	/*
	 * src string to be encrypted
	 * key Key used for encryption Key length can be any one of 128bit, 192bit, 256bit
	 * 16-bit key corresponds to 128bit
	 */
	key := "@TRUE-IT^!#$+*&%@TRUE-IT^!#$+*&%"
	data := "<Paravit#123/ทำการทดสอบ!!>"

	encryptedText, err := cryptography.EncryptAESECB(key, data)

	if err != nil {
		t.Errorf("TestAES256Ecb - EncryptAESECB Error %v", err.Error())
	} else {
		fmt.Println("encryptedText:", encryptedText)
	}

	plainText, err := cryptography.DecryptAESECB(key, encryptedText)

	if err != nil {
		t.Errorf("TestAES256Ecb - DecryptAESECB Error %v", err.Error())
	} else if data != plainText {
		t.Errorf("data not equal plainText")
	} else {
		fmt.Println("plantext:", plainText)
	}
}

func TestAES256Cbc(t *testing.T) {
	key := "2b7e276890cfg2a6abf7158809cf4f3c"
	data := "<Paravit#123/ทำการทดสอบ!!>"
	iv := []byte("0000000000000000")
	encryptedText, err := cryptography.EncryptAESCBC(key, data, iv)

	if err != nil {
		t.Errorf("TestAES256Cbc - EncryptAESCBC Error %v", err.Error())
	} else {
		fmt.Println("encryptedText:", encryptedText)
	}

	plainText, err := cryptography.DecryptAESCBC(key, encryptedText, iv)

	if err != nil {
		t.Errorf("TestAES256Ecb - DecryptAESCBC Error %v", err.Error())
	} else if data != plainText {
		t.Errorf("data not equal plainText")
	} else {
		fmt.Println("plainText:", plainText)
	}
}

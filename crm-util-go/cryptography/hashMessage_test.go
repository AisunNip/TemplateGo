package cryptography_test

import (
	"crm-util-go/cryptography"
	"fmt"
	"testing"
)

func TestHashMessage(t *testing.T) {
	secretKey := "#True-IT!!2022**"
	salt := "ABCDEF1234567890xxx"
	data := "<Paravit#123/ทำการทดสอบ!!>"

	fmt.Println("MD5:", cryptography.MD5(salt, data))
	fmt.Println("SHA1:", cryptography.SHA1(salt, data))
	fmt.Println("SHA256:", cryptography.SHA256(salt, data))
	fmt.Println("SHA384:", cryptography.SHA384(salt, data))
	fmt.Println("SHA512:", cryptography.SHA512(salt, data))
	fmt.Println("SHA512/224:", cryptography.SHA512_224(salt, data))
	fmt.Println("SHA512/256:", cryptography.SHA512_256(salt, data))
	fmt.Println("HMACSHA512:", cryptography.HMACSHA512(secretKey, data))
}

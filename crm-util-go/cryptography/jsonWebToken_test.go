package cryptography_test

import (
	"crm-util-go/crmencoding"
	"crm-util-go/cryptography"
	"fmt"
	"testing"
	"time"
)

func getInitJwtClaims() cryptography.JwtClaims {
	var jwtClaims cryptography.JwtClaims

	jwtClaims.JwtID = "123jwtid"
	jwtClaims.Subject = "test subject"

	// Test Expire
	//jwtClaims.ExpirationTime = time.Now().AddDate(0, 0, -1)
	jwtClaims.ExpirationTime = time.Now().Add(time.Minute * 15)

	jwtClaims.UserID = "12345"
	jwtClaims.FirstName = "Paravit"
	jwtClaims.LastName = "Tunvichian"
	jwtClaims.MiddleName = "Pui"
	jwtClaims.Email = "paravit_tun@truecorp.co.th"

	return jwtClaims
}

func TestCreateJWTokenHS512(t *testing.T) {
	jwtClaims := getInitJwtClaims()
	secret := "abc56789"

	fmt.Println("########## CreateJWTokenHS512 ##########")
	//token, err := cryptography.CreateJWTokenHS256(secret, jwtClaims)
	//token, err := cryptography.CreateJWTokenHS384(secret, jwtClaims)
	token, err := cryptography.CreateJWTokenHS512(secret, jwtClaims)

	if err != nil {
		t.Errorf("VerifyToken Error: " + err.Error())
	} else {
		fmt.Println("token: " + token)
	}

	fmt.Println("########## VerifyJWToken ##########")
	jwtClaimsOK, err := cryptography.VerifyJWToken(secret, token)

	if err != nil {
		t.Errorf("VerifyToken Error: " + err.Error())
	} else {
		jsonString, _ := crmencoding.StructToJson(jwtClaimsOK)
		fmt.Println("JWT:", jsonString)
		fmt.Println("ExpDate", jwtClaimsOK.ExpirationTime, jwtClaimsOK.ExpirationTime.Unix())
	}
	fmt.Println("##################################")

}

func TestCreateJWTokenRS512(t *testing.T) {
	fmt.Println("########## CreateJWTokenRS512 ##########")
	jwtClaims := getInitJwtClaims()
	tokenRS512, err := cryptography.CreateJWTokenRS512("../certfile/mycompany.key-pkcs8.pem", jwtClaims)

	if err != nil {
		t.Errorf("CreateJWTokenRS512 Error: " + err.Error())
	} else {
		fmt.Println("RS512 token: " + tokenRS512)
	}

	fmt.Println("########## VerifyJWTokenRS512 ##########")
	jwtClaimsOKRS512, err := cryptography.VerifyJWTokenRS512("../certfile/mycompany.pubkey.pem", tokenRS512)

	if err != nil {
		t.Errorf("VerifyJWTokenRS512 Error: " + err.Error())
	} else {
		fmt.Println("VerifyJWTokenRS512 success")
		fmt.Println("UserID: " + jwtClaimsOKRS512.UserID)
		fmt.Println("FirstName: " + jwtClaimsOKRS512.FirstName)
		fmt.Println("LastName: " + jwtClaimsOKRS512.LastName)
		fmt.Println("Email: " + jwtClaimsOKRS512.Email)
	}
}

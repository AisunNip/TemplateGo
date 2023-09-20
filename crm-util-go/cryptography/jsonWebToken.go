package cryptography

import (
	"crm-util-go/validate"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strconv"
	"time"
)

type JwtClaims struct {
	JwtID          string
	Issuer         string
	IssuedTime     time.Time
	Subject        string
	NotBeforeTime  time.Time
	ExpirationTime time.Time
	Audience       string
	UserID         string
	FullName       string
	FirstName      string
	LastName       string
	MiddleName     string
	Gender         string
	Email          string
}

func jwtClaimsToMap(jwtClaims JwtClaims) jwt.MapClaims {
	claims := jwt.MapClaims{}

	if validate.HasStringValue(jwtClaims.JwtID) {
		claims["jti"] = jwtClaims.JwtID
	}

	if validate.HasStringValue(jwtClaims.Issuer) {
		claims["iss"] = jwtClaims.Issuer
	}

	if validate.HasDateTime(jwtClaims.IssuedTime) {
		claims["iat"] = jwtClaims.IssuedTime.Unix()
	}

	if validate.HasStringValue(jwtClaims.Subject) {
		claims["sub"] = jwtClaims.Subject
	}

	if validate.HasDateTime(jwtClaims.NotBeforeTime) {
		claims["nbf"] = jwtClaims.NotBeforeTime.Unix()
	}

	if validate.HasDateTime(jwtClaims.ExpirationTime) {
		claims["exp"] = jwtClaims.ExpirationTime.Unix()
	}

	if validate.HasStringValue(jwtClaims.Audience) {
		claims["aud"] = jwtClaims.Audience
	}

	if validate.HasStringValue(jwtClaims.UserID) {
		claims["user_id"] = jwtClaims.UserID
	}

	if validate.HasStringValue(jwtClaims.FullName) {
		claims["name"] = jwtClaims.FullName
	}

	if validate.HasStringValue(jwtClaims.FirstName) {
		claims["given_name"] = jwtClaims.FirstName
	}

	if validate.HasStringValue(jwtClaims.LastName) {
		claims["family_name"] = jwtClaims.LastName
	}

	if validate.HasStringValue(jwtClaims.MiddleName) {
		claims["middle_name"] = jwtClaims.MiddleName
	}

	if validate.HasStringValue(jwtClaims.Gender) {
		claims["gender"] = jwtClaims.Gender
	}

	if validate.HasStringValue(jwtClaims.Email) {
		claims["email"] = jwtClaims.Email
	}

	return claims
}

func mapToJwtClaims(mapClaims jwt.MapClaims) JwtClaims {
	var jwtClaims JwtClaims

	for key, val := range mapClaims {
		switch key {
		case "jti":
			jwtClaims.JwtID = fmt.Sprintf("%v", val)
		case "iss":
			jwtClaims.Issuer = fmt.Sprintf("%v", val)
		case "iat":
			jwtClaims.IssuedTime = interfaceToTimeUnix(val)
		case "sub":
			jwtClaims.Subject = fmt.Sprintf("%v", val)
		case "nbf":
			jwtClaims.NotBeforeTime = interfaceToTimeUnix(val)
		case "exp":
			jwtClaims.ExpirationTime = interfaceToTimeUnix(val)
		case "aud":
			jwtClaims.Audience = fmt.Sprintf("%v", val)
		case "user_id":
			jwtClaims.UserID = fmt.Sprintf("%v", val)
		case "name":
			jwtClaims.FullName = fmt.Sprintf("%v", val)
		case "given_name":
			jwtClaims.FirstName = fmt.Sprintf("%v", val)
		case "family_name":
			jwtClaims.LastName = fmt.Sprintf("%v", val)
		case "middle_name":
			jwtClaims.MiddleName = fmt.Sprintf("%v", val)
		case "gender":
			jwtClaims.Gender = fmt.Sprintf("%v", val)
		case "email":
			jwtClaims.Email = fmt.Sprintf("%v", val)
		default:
			fmt.Sprintf("%v", val)
		}
	}

	return jwtClaims
}

func interfaceToTimeUnix(data interface{}) time.Time {
	var result time.Time

	if data != nil {
		dataDT, err := strconv.ParseFloat(fmt.Sprintf("%g", data), 64)

		if err == nil {
			result = time.Unix(int64(dataDT), 0)
		}
	}

	return result
}

func createJWToken(signingMethod jwt.SigningMethod, secret string, jwtClaims JwtClaims) (string, error) {
	claims := jwtClaimsToMap(jwtClaims)
	jwtToken := jwt.NewWithClaims(signingMethod, claims)
	token, err := jwtToken.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func CreateJWTokenHS256(secret string, jwtClaims JwtClaims) (string, error) {
	return createJWToken(jwt.SigningMethodHS256, secret, jwtClaims)
}

func CreateJWTokenHS384(secret string, jwtClaims JwtClaims) (string, error) {
	return createJWToken(jwt.SigningMethodHS384, secret, jwtClaims)
}

func CreateJWTokenHS512(secret string, jwtClaims JwtClaims) (string, error) {
	return createJWToken(jwt.SigningMethodHS512, secret, jwtClaims)
}

func CreateJWTokenRS512(privateKeyPEM string, jwtClaims JwtClaims) (string, error) {
	/*
		1. Create a self signed certificate
			openssl req -x509 -sha512 -newkey rsa:4096 \
		   		-keyout mycompany.key.pem -out mycompany.cert.pem -days 365
		2. Convert the private key to PKCS#8 format.
			openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in mycompany.key.pem -out mycompany.key-pkcs8.pem
	*/
	privateKey, err := os.ReadFile(privateKeyPEM)

	if err != nil {
		fmt.Println("Unable to read RSA private key: " + err.Error())
		return "", err
	}

	var rsaKey *rsa.PrivateKey
	rsaKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	if err != nil {
		fmt.Println("Unable to parse RSA private key: " + err.Error())
		return "", err
	}

	claims := jwtClaimsToMap(jwtClaims)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	token, err := jwtToken.SignedString(rsaKey)

	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyJWToken(secret string, tokenString string) (JwtClaims, error) {
	var jwtClaims JwtClaims

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err == nil {
		jwtClaims = mapToJwtClaims(claims)
	}

	return jwtClaims, err
}

func VerifyJWTokenRS512(publicKeyPEM string, tokenString string) (JwtClaims, error) {
	var jwtClaims JwtClaims

	publicKey, err := os.ReadFile(publicKeyPEM)

	if err != nil {
		fmt.Println("Unable to read RSA public key: " + err.Error())
		return jwtClaims, err
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(publicKey)
	})

	if err == nil {
		jwtClaims = mapToJwtClaims(claims)
	}

	return jwtClaims, err
}

package cryptography

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
)

func hashMessage(hasher hash.Hash, salt string, msg string) string {
	if len(salt) > 0 {
		io.WriteString(hasher, salt)
	}

	io.WriteString(hasher, msg)

	return hex.EncodeToString(hasher.Sum(nil))
}

/*
	Message digests (MD) are secure one-way hash functions.
*/
func MD5(salt string, msg string) string {
	return hashMessage(md5.New(), salt, msg)
}

func SHA1(salt string, msg string) string {
	return hashMessage(sha1.New(), salt, msg)
}

func SHA256(salt string, msg string) string {
	return hashMessage(sha256.New(), salt, msg)
}

func SHA384(salt string, msg string) string {
	return hashMessage(sha512.New384(), salt, msg)
}

func SHA512(salt string, msg string) string {
	return hashMessage(sha512.New(), salt, msg)
}

func SHA512_224(salt string, msg string) string {
	return hashMessage(sha512.New512_224(), salt, msg)
}

func SHA512_256(salt string, msg string) string {
	return hashMessage(sha512.New512_256(), salt, msg)
}

func HMACSHA512(secretKey string, msg string) string {
	hasher := hmac.New(sha512.New, []byte(secretKey))
	return hashMessage(hasher, "", msg)
}

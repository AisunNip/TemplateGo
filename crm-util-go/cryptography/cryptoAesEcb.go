package cryptography

/*
   cryptoAesECB.go implement AES/ECB/PKCS5Padding same Java language
   Key used for encryption Key length can be any one of 128bit, 192bit, 256bit
   16-bit key corresponds to 128bit
*/
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func EncryptAESECB(key string, plainText string) (encryptedText string, err error) {
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return encryptedText, err
	}

	if plainText == "" {
		return encryptedText, err
	}

	ecbBlockMode := NewECBEncrypter(block)
	binaryPlainText := []byte(plainText)
	binaryPlainText = PKCS5Padding(binaryPlainText, block.BlockSize())
	encrypted := make([]byte, len(binaryPlainText))
	ecbBlockMode.CryptBlocks(encrypted, binaryPlainText)

	encryptedText = base64.StdEncoding.EncodeToString(encrypted)

	return encryptedText, nil
}

func DecryptAESECB(key string, encryptedText string) (plainText string, err error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedText)

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return plainText, err
	}

	ecbBlockMode := NewECBDecrypter(block)
	binaryPlainText := make([]byte, len(encrypted))
	ecbBlockMode.CryptBlocks(binaryPlainText, encrypted)
	binaryPlainText = PKCS5Unpadding(binaryPlainText)

	plainText = string(binaryPlainText)

	return plainText, err
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	// remove the last byte unpadding times
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}

	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}

	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

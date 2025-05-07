package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"text-database/pkg/utilities"
)

type SecureTextEncoder struct {
	key []byte
}

var globalEncoderKey SecureTextEncoder

func NewSecureTextEncoder(secretKey string) *SecureTextEncoder {
	hasher := sha256.New()
	hasher.Write([]byte(secretKey))
	key := hasher.Sum(nil)

	return &SecureTextEncoder{
		key: key,
	}
}

func (e *SecureTextEncoder) Encode(plainText string) (string, error) {
	// Create cipher block
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure
	// Use GCM mode for authenticated encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Never use more than 2^32 random nonces with a given key
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and seal the data
	ciphertext := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	// Convert to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *SecureTextEncoder) Decode(encodedText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
func (e *SecureTextEncoder) ReadAndDecode(dbName string) string {
	data := utilities.Must(os.ReadFile(dbName))
	encodeData := utilities.Must(e.Decode(string(data)))
	return encodeData
}

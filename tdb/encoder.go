package tdb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strings"
)

type secureTextEncoder struct {
	key []byte
}

var globalEncoderKey secureTextEncoder

// newSecureTextEncoder creates a new secureTextEncoder instance with the provided secret key.
// The secret key is hashed using SHA-256 to create the encryption key.
func newSecureTextEncoder(secretKey string) *secureTextEncoder {
	hasher := sha256.New()
	hasher.Write([]byte(secretKey))
	key := hasher.Sum(nil)

	return &secureTextEncoder{
		key: key,
	}
}

// Encode encrypts the plain text using AES-GCM encryption and returns the base64-encoded result
// with "ENG" prefix. Returns error if encryption fails.
func (e *secureTextEncoder) Encode(plainText string) (string, error) {
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
	b64 := base64.StdEncoding.EncodeToString(ciphertext)
	//Add Prefix
	prefix := "ENG" + b64
	return prefix, nil

}

// Decode decrypts the encoded text (with "ENG" prefix removed) using AES-GCM decryption.
// Returns the original plain text or error if decryption fails.
func (e *secureTextEncoder) Decode(encodedText string) (string, error) {
	encodedText = strings.Replace(encodedText, "ENG", "", 1)
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

// readAndDecode reads the content of the database file and decodes it if encryption is enabled.
// Returns the decoded content as a string.
func (e *secureTextEncoder) readAndDecode(dbName string) string {
	data := string(must(os.ReadFile(dbName)))
	if encryptionKeyExist {
		data = must(e.Decode(data))
	}
	return data
}

// isEncode checks if the given text is encoded by verifying if it starts with "ENG" prefix.
// Returns true if the text is encoded, false otherwise.
func isEncode(text string) bool {
	if strings.HasPrefix(text, "ENG") {
		return true
	}
	return false
}

// encodeAndSave encrypts the provided data using the global encoder key and saves it to the database file.
// Panics if encryption or file writing fails.
func encodeAndSave(data string) {
	encodeData := must(globalEncoderKey.Encode(data))
	errorHandler(os.WriteFile(dbName, []byte(encodeData), 0644))
}

// decodeAndSave decrypts the provided data using the global encoder key and saves it to the database file.
// Panics if decryption or file writing fails.
func decodeAndSave(data string) {
	decodeData := must(globalEncoderKey.Decode(data))
	errorHandler(os.WriteFile(dbName, []byte(decodeData), 0644))
}

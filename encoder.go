package text_database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/sheymor21/tdb/utilities"
	"io"
	"os"
	"strings"
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
	b64 := base64.StdEncoding.EncodeToString(ciphertext)
	//Add Prefix
	prefix := "ENG" + b64
	return prefix, nil

}

func (e *SecureTextEncoder) Decode(encodedText string) (string, error) {
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
func (e *SecureTextEncoder) ReadAndDecode(dbName string) string {
	data := string(utilities.Must(os.ReadFile(dbName)))
	if encryptionKeyExist {
		data = utilities.Must(e.Decode(data))
	}
	return data
}

func IsEncode(text string) bool {
	if strings.HasPrefix(text, "ENG") {
		return true
	}
	return false
}

func EncodeAndSave(data string) {
	encodeData := utilities.Must(globalEncoderKey.Encode(data))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(encodeData), 0644))
}

func DecodeAndSave(data string) {
	decodeData := utilities.Must(globalEncoderKey.Decode(data))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(decodeData), 0644))
}

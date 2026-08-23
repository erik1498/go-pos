package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"go-pos/internal/domain"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	cipherFormat = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
)

func GenerateBlindedIndex(data, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func EncryptAES(plaintext string, secretKey []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, nil
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	return cipherText, nil
}

func DecryptAES(cipherText []byte, secretKey []byte) (string, error) {
	if len(cipherText) == 0 {
		return "", nil
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return "", domain.ErrChiperTextToShort
	}

	nonce, cipherTextMessage := cipherText[:nonceSize], cipherText[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherTextMessage, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

type argonParams struct {
	memory     uint32
	time       uint32
	threads    uint8
	saltLength uint32
	keyLength  uint32
}

func getEnvAsUint32(key string, fallback uint32) uint32 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	valueInt, err := strconv.ParseUint(valueStr, 10, 32)
	if err != nil {
		return fallback
	}
	return uint32(valueInt)
}

func getArgonParams() argonParams {
	return argonParams{
		memory:     getEnvAsUint32("ARGON2_MEMORY", 64*1024),
		time:       getEnvAsUint32("ARGON2_TIME", 3),
		threads:    uint8(getEnvAsUint32("ARGON2_THREADS", 4)),
		saltLength: getEnvAsUint32("ARGON2_SALT_LENGTH", 16),
		keyLength:  getEnvAsUint32("ARGON2_KEY_LENGTH", 32),
	}
}

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)

	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	return salt, nil
}

func HashPassword(password string, aesKey []byte) (string, error) {
	p := getArgonParams()
	salt, err := generateSalt(p.saltLength)
	if err != nil {
		return "", err
	}

	argonHash := argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threads, p.keyLength)

	argonHashEncrypt, err := EncryptAES(string(argonHash), aesKey)
	if err != nil {
		return "", err
	}
	b64Hash := base64.RawStdEncoding.EncodeToString(argonHashEncrypt)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	encodedHash := fmt.Sprintf(cipherFormat, argon2.Version, p.memory, p.time, p.threads, b64Hash, b64Salt)

	return encodedHash, nil
}

func VerifyPassword(password string, hash string, aesKey []byte) (bool, error) {
	parts := strings.Split(hash, "$")

	var memory, time uint32
	var threads uint8

	switch parts[1] {
	case "argon2id":
		_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
		if err != nil {
			return false, err
		}

		salt, err := base64.RawStdEncoding.DecodeString(parts[5])
		if err != nil {
			return false, err
		}

		hash, err := base64.RawStdEncoding.DecodeString(parts[4])
		if err != nil {
			return false, err
		}

		decryptHash, err := DecryptAES(hash, aesKey)
		if err != nil {
			return false, err
		}

		keyLen := uint32(len(decryptHash))

		comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

		return subtle.ConstantTimeCompare(comparisonHash, []byte(decryptHash)) == 1, nil
	}
	return true, nil
}

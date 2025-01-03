package license_handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
)

type LicenseHandler struct {
	filePath      string
	encryptionKey []byte // 32 bytes for AES-256
}

func NewLicenseHandler(path string) *LicenseHandler {
	key := make([]byte, 32)
	// Generate static key - in production use secure key management
	copy(key, []byte("blizzflow-static-encryption-key-32b"))

	return &LicenseHandler{
		filePath:      path,
		encryptionKey: key,
	}
}

func (h *LicenseHandler) SaveLicense(licenseKey string) error {
	// Encrypt
	block, err := aes.NewCipher(h.encryptionKey)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Encrypt and encode
	ciphertext := gcm.Seal(nonce, nonce, []byte(licenseKey), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	// Save to file
	return os.WriteFile(h.filePath, []byte(encoded), 0644)
}

func (h *LicenseHandler) ReadLicense() (string, error) {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(h.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

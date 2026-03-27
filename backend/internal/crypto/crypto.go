package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

type CryptoService struct {
	keySize         int
	argon2Time      uint32
	argon2Mem       uint32
	argon2Threads   uint8
}

func NewCryptoService() *CryptoService {
	return &CryptoService{
		keySize:       32,
		argon2Time:    3,
		argon2Mem:     65536,
		argon2Threads: 4,
	}
}

func (s *CryptoService) DeriveKey(password, salt []byte) ([]byte, error) {
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}
	if len(salt) != 16 {
		return nil, fmt.Errorf("salt must be 16 bytes")
	}

	key := argon2.IDKey(password, salt, s.argon2Time, s.argon2Mem, s.argon2Threads, uint32(s.keySize))
	return key, nil
}

func (s *CryptoService) Encrypt(plaintext, key []byte) ([]byte, error) {
	if len(key) != s.keySize {
		return nil, fmt.Errorf("key must be %d bytes", s.keySize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (s *CryptoService) Decrypt(ciphertext, key []byte) ([]byte, error) {
	if len(key) != s.keySize {
		return nil, fmt.Errorf("key must be %d bytes", s.keySize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func (s *CryptoService) EncryptToB64(plaintext, key []byte) (string, error) {
	ciphertext, err := s.Encrypt(plaintext, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *CryptoService) DecryptFromB64(ciphertextB64 string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return s.Decrypt(ciphertext, key)
}

func (s *CryptoService) GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func (s *CryptoService) Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func (s *CryptoService) GetKeySize() int {
	return s.keySize
}

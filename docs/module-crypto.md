# Crypto Module (Go)

## 1. Purpose

Provides cryptographic utilities for secure password hashing and data encryption. Belongs to the **Infrastructure layer**.

## 2. Scope

**Functions provided:**
- Key derivation (Argon2id)
- Symmetric encryption (AES-256-GCM)
- Hashing (SHA-256)
- Random generation

## 3. Public API

```go
type CryptoService struct{}

func (s *CryptoService) DeriveKey(password, salt []byte) ([]byte, error)
func (s *CryptoService) Encrypt(plaintext, key []byte) ([]byte, error)
func (s *CryptoService) Decrypt(ciphertext, key []byte) ([]byte, error)
func (s *CryptoService) EncryptToB64(plaintext, key []byte) (string, error)
func (s *CryptoService) DecryptFromB64(ciphertextB64 string, key []byte) ([]byte, error)
func (s *CryptoService) GenerateSalt() ([]byte, error)
func (s *CryptoService) Hash(data []byte) []byte
func (s *CryptoService) GetKeySize() int
```

## 4. Implementation Details

### Key Derivation

Uses Argon2id algorithm:
- **Memory**: 64 MB
- **Iterations**: 3
- **Parallelism**: 4 threads
- **Output**: 32 bytes (256-bit key)

### Symmetric Encryption

Uses AES-256-GCM:
- **Key size**: 32 bytes (256-bit)
- **Nonce**: 12 bytes (random per encryption)
- **Authentication**: GCM mode provides authenticity

### Salt Generation

16 bytes of cryptographically secure random data.

## 5. Security Considerations

- **Timing attacks**: Constant-time comparison not implemented (use subtle.ConstantTimeCompare)
- **Key stretching**: Argon2id resistant to GPU/ASIC attacks
- **Nonce reuse prevention**: Random nonce per encryption
- **No custom crypto**: Uses standard Go crypto library

## 6. Integration

**Used by:**
- Auth module (password hashing)
- Vault module (entry encryption)

## 7. Roadmap

- [ ] Add constant-time comparison
- [ ] Key rotation support
- [ ] Asymmetric encryption for sharing
- [ ] Hardware security module integration

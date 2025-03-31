package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

func GenerateSalt() ([]byte, error) {
	length := 16
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKey gera uma chave a partir da senha e do salt usando scrypt.
func DeriveKey(password, salt []byte) ([]byte, error) {
	// Parâmetros: N=32768, r=8, p=1, chave de 32 bytes (AES-256)
	return scrypt.Key(password, salt, 32768, 8, 1, 32)
}

// Encrypt criptografa o plaintext usando AES-GCM.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return []byte{}, nil
	}

	block, err := aes.NewCipher(key)
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
	// O nonce é prefixado ao ciphertext
	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt descriptografa o ciphertext usando AES-GCM.
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return []byte{}, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext muito curto")
	}
	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return aesGCM.Open(nil, nonce, ct, nil)
}

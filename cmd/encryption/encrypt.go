package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
)

type EncryptedMessage struct {
	Nonce     []byte `json:"Nonce"`
	PlainText []byte `json:"PlainText"`
	Key       []byte `json:"Key"`
}

func GCMEncrypter(key []byte, message []byte) ([]byte, []byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, message, nil)
	return ciphertext, nonce, nil
}

func EncryptHybrid(plaintext []byte, pub *rsa.PublicKey) ([]byte, error) {
	var encryptedMessage EncryptedMessage
	// generate random AES-256 key
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, fmt.Errorf(" generate aes key: %w", err)
	}

	// encrypt AES key with RSA-OAEP (SHA-256)
	hash := sha256.New()
	encKey, err := rsa.EncryptOAEP(hash, rand.Reader, pub, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("encrypt aes key: %w", err)
	}

	// encrypt plaintext with AES-GCM
	ciphertext, nonce, err := GCMEncrypter(aesKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("encrypt aes key: %w", err)
	}

	// marshal json response
	encryptedMessage.Key = encKey
	encryptedMessage.PlainText = ciphertext
	encryptedMessage.Nonce = nonce

	jsonMessage, err := json.Marshal(encryptedMessage)
	if err != nil {
		return nil, fmt.Errorf("failed marshal encrypted message: %w", err)
	}
	return jsonMessage, nil
}

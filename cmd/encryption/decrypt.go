package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
)

func GCMDecrypter(ciphertext []byte, nonce []byte, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func DecryptHybrid(key *rsa.PrivateKey, message []byte) ([]byte, error) {
	var encryptedMessage EncryptedMessage
	if err := json.Unmarshal(message, &encryptedMessage); err != nil {
		return nil, err
	}

	// decrypt symmetric key by secret key
	aeskey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, encryptedMessage.Key, nil)
	if err != nil {
		return nil, err
	}

	// decrypt message data by decrypted symmetric key
	plainText, err := GCMDecrypter(encryptedMessage.PlainText, encryptedMessage.Nonce, aeskey)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

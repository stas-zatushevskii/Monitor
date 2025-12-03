package encryption

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestGCMEncrypterAndDecrypter(t *testing.T) {
	plaintext := []byte("Hello, world! This is a test message.")

	// generate AES-256 key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	// encrypt
	ciphertext, nonce, err := GCMEncrypter(key, plaintext)
	if err != nil {
		t.Fatalf("GCMEncrypter failed: %v", err)
	}

	// decrypt
	decrypted, err := GCMDecrypter(ciphertext, nonce, key)
	if err != nil {
		t.Fatalf("GCMDecrypter failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text does not match original.\nExpected: %s\nGot: %s", plaintext, decrypted)
	}

	// try decrypt with wrong key
	wrongKey := make([]byte, 32)
	rand.Read(wrongKey)
	_, err = GCMDecrypter(ciphertext, nonce, wrongKey)
	if err == nil {
		t.Errorf("Decryption with wrong key should fail but succeeded")
	}
}

func TestHybridEncryptDecrypt(t *testing.T) {
	// generate RSA key pair
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA keys: %v", err)
	}
	pubKey := &privKey.PublicKey

	plaintext := []byte("Hybrid encryption test message.")

	// encrypt hybrid
	encrypted, err := EncryptHybrid(plaintext, pubKey)
	if err != nil {
		t.Fatalf("EncryptHybrid failed: %v", err)
	}

	// decrypt hybrid
	decrypted, err := DecryptHybrid(privKey, encrypted)
	if err != nil {
		t.Fatalf("DecryptHybrid failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Hybrid decrypted text does not match original.\nExpected: %s\nGot: %s", plaintext, decrypted)
	}
}

func TestHybridTampering(t *testing.T) {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &privKey.PublicKey

	plaintext := []byte("Message for tampering test")
	encrypted, _ := EncryptHybrid(plaintext, pubKey)

	// tamper with the ciphertext JSON
	encrypted[50] ^= 0xFF // flip a byte

	_, err := DecryptHybrid(privKey, encrypted)
	if err == nil {
		t.Errorf("DecryptHybrid should fail for tampered message but succeeded")
	}
}

func TestMultipleMessages(t *testing.T) {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &privKey.PublicKey

	messages := [][]byte{
		[]byte("Short msg"),
		[]byte("Another test message"),
		[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
	}

	for _, msg := range messages {
		enc, err := EncryptHybrid(msg, pubKey)
		if err != nil {
			t.Fatalf("EncryptHybrid failed: %v", err)
		}
		dec, err := DecryptHybrid(privKey, enc)
		if err != nil {
			t.Fatalf("DecryptHybrid failed: %v", err)
		}
		if !bytes.Equal(msg, dec) {
			t.Errorf("Mismatch\nExpected: %s\nGot: %s", msg, dec)
		}
	}
}

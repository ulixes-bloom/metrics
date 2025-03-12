// Package rsa privides functions to encrypt end decrypt text using RSA algorythm.
// It also allows to generate public-private key pair
package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// Encrypt encrypts plaintext using public key cpecified in publicKeyPath
func Encrypt(plaintext []byte, publicKeyPath string) ([]byte, error) {
	publicKeyPem, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("cipher.encrypt.readFile '%s': %w", publicKeyPath, err)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPem)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("cipher.encrypt.parsePublicKey '%s': %w", publicKeyPath, err)
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), plaintext)
	if err != nil {
		return nil, fmt.Errorf("cipher.encrypt: %w", err)
	}

	return ciphertext, nil
}

// Decrypt decrypts ciphertest using private key cpecified in privateKeyPath
func Decrypt(ciphertext []byte, privateKeyPath string) ([]byte, error) {
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("cipher.decrypt.readFile '%s': %w", privateKeyPath, err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("cipher.decrypt.parsePrivateKey '%s': %w", privateKeyPath, err)
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("cipher.decrypt: %w", err)
	}

	return plaintext, nil
}

// GenerateKeyPair generates random RSA public and private keys in pem format
func GenerateKeyPair() (publicKeyPEM, privateKeyPEM []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	publicKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return
}

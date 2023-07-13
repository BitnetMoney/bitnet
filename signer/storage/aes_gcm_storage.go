// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/log"
)

type storedCredential struct {
	// The iv
	Iv []byte `json:"iv"`
	// The ciphertext
	CipherText []byte `json:"c"`
}

// AESEncryptedStorage is a storage type which is backed by a json-file. The json-file contains
// key-value mappings, where the keys are _not_ encrypted, only the values are.
type AESEncryptedStorage struct {
	// File to read/write credentials
	filename string
	// Key stored in base64
	key []byte
}

// NewAESEncryptedStorage creates a new encrypted storage backed by the given file/key
func NewAESEncryptedStorage(filename string, key []byte) *AESEncryptedStorage {
	return &AESEncryptedStorage{
		filename: filename,
		key:      key,
	}
}

// Put stores a value by key. 0-length keys results in noop.
func (s *AESEncryptedStorage) Put(key, value string) {
	if len(key) == 0 {
		return
	}
	data, err := s.readEncryptedStorage()
	if err != nil {
		log.Warn("Failed to read encrypted storage", "err", err, "file", s.filename)
		return
	}
	ciphertext, iv, err := encrypt(s.key, []byte(value), []byte(key))
	if err != nil {
		log.Warn("Failed to encrypt entry", "err", err)
		return
	}
	encrypted := storedCredential{Iv: iv, CipherText: ciphertext}
	data[key] = encrypted
	if err = s.writeEncryptedStorage(data); err != nil {
		log.Warn("Failed to write entry", "err", err)
	}
}

// Get returns the previously stored value, or an error if it does not exist or
// key is of 0-length.
func (s *AESEncryptedStorage) Get(key string) (string, error) {
	if len(key) == 0 {
		return "", ErrZeroKey
	}
	data, err := s.readEncryptedStorage()
	if err != nil {
		log.Warn("Failed to read encrypted storage", "err", err, "file", s.filename)
		return "", err
	}
	encrypted, exist := data[key]
	if !exist {
		log.Warn("Key does not exist", "key", key)
		return "", ErrNotFound
	}
	entry, err := decrypt(s.key, encrypted.Iv, encrypted.CipherText, []byte(key))
	if err != nil {
		log.Warn("Failed to decrypt key", "key", key)
		return "", err
	}
	return string(entry), nil
}

// Del removes a key-value pair. If the key doesn't exist, the method is a noop.
func (s *AESEncryptedStorage) Del(key string) {
	data, err := s.readEncryptedStorage()
	if err != nil {
		log.Warn("Failed to read encrypted storage", "err", err, "file", s.filename)
		return
	}
	delete(data, key)
	if err = s.writeEncryptedStorage(data); err != nil {
		log.Warn("Failed to write entry", "err", err)
	}
}

// readEncryptedStorage reads the file with encrypted creds
func (s *AESEncryptedStorage) readEncryptedStorage() (map[string]storedCredential, error) {
	creds := make(map[string]storedCredential)
	raw, err := os.ReadFile(s.filename)

	if err != nil {
		if os.IsNotExist(err) {
			// Doesn't exist yet
			return creds, nil
		}
		log.Warn("Failed to read encrypted storage", "err", err, "file", s.filename)
	}
	if err = json.Unmarshal(raw, &creds); err != nil {
		log.Warn("Failed to unmarshal encrypted storage", "err", err, "file", s.filename)
		return nil, err
	}
	return creds, nil
}

// writeEncryptedStorage write the file with encrypted creds
func (s *AESEncryptedStorage) writeEncryptedStorage(creds map[string]storedCredential) error {
	raw, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	if err = os.WriteFile(s.filename, raw, 0600); err != nil {
		return err
	}
	return nil
}

// encrypt encrypts plaintext with the given key, with additional data
// The 'additionalData' is used to place the (plaintext) KV-store key into the V,
// to prevent the possibility to alter a K, or swap two entries in the KV store with each other.
func encrypt(key []byte, plaintext []byte, additionalData []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, additionalData)
	return ciphertext, nonce, nil
}

func decrypt(key []byte, nonce []byte, ciphertext []byte, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

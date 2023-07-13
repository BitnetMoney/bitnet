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

import "errors"

var (
	// ErrZeroKey is returned if an attempt was made to inset a 0-length key.
	ErrZeroKey = errors.New("0-length key")

	// ErrNotFound is returned if an unknown key is attempted to be retrieved.
	ErrNotFound = errors.New("not found")
)

type Storage interface {
	// Put stores a value by key. 0-length keys results in noop.
	Put(key, value string)

	// Get returns the previously stored value, or an error if the key is 0-length
	// or unknown.
	Get(key string) (string, error)

	// Del removes a key-value pair. If the key doesn't exist, the method is a noop.
	Del(key string)
}

// EphemeralStorage is an in-memory storage that does
// not persist values to disk. Mainly used for testing
type EphemeralStorage struct {
	data map[string]string
}

// Put stores a value by key. 0-length keys results in noop.
func (s *EphemeralStorage) Put(key, value string) {
	if len(key) == 0 {
		return
	}
	s.data[key] = value
}

// Get returns the previously stored value, or an error if the key is 0-length
// or unknown.
func (s *EphemeralStorage) Get(key string) (string, error) {
	if len(key) == 0 {
		return "", ErrZeroKey
	}
	if v, ok := s.data[key]; ok {
		return v, nil
	}
	return "", ErrNotFound
}

// Del removes a key-value pair. If the key doesn't exist, the method is a noop.
func (s *EphemeralStorage) Del(key string) {
	delete(s.data, key)
}

func NewEphemeralStorage() Storage {
	s := &EphemeralStorage{
		data: make(map[string]string),
	}
	return s
}

// NoStorage is a dummy construct which doesn't remember anything you tell it
type NoStorage struct{}

func (s *NoStorage) Put(key, value string) {}
func (s *NoStorage) Del(key string)        {}
func (s *NoStorage) Get(key string) (string, error) {
	return "", errors.New("missing key, I probably forgot")
}

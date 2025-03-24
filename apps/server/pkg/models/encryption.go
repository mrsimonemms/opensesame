/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5" //nolint:gosec // use md5 to create 32 character string
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

func decrypt(input, key string) ([]byte, error) {
	gcmInstance, err := createCipher(key)
	if err != nil {
		return nil, err
	}
	nonceSize := gcmInstance.NonceSize()

	ciphered, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, fmt.Errorf("error decoding cipher from base64: %w", err)
	}

	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return nil, fmt.Errorf("error deciphering: %w", err)
	}

	return originalText, nil
}

func createCipher(key string) (cipher.AEAD, error) {
	aesBlock, err := aes.NewCipher([]byte(hash(key)))
	if err != nil {
		return nil, fmt.Errorf("error creating aes cipher block: %w", err)
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, fmt.Errorf("error creating gcm block: %w", err)
	}

	return gcmInstance, nil
}

func encrypt(input, key string) ([]byte, error) {
	gcmInstance, err := createCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("error reading string: %w", err)
	}

	sealedCipher := gcmInstance.Seal(nonce, nonce, []byte(input), nil)

	return []byte(base64.StdEncoding.EncodeToString(sealedCipher)), nil
}

func hash(input string) string {
	plainText := []byte(input)
	//nolint:gosec // use md5 to create 32 character string
	md5Hash := md5.Sum(plainText)
	return hex.EncodeToString(md5Hash[:])
}

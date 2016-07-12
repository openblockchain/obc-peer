/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package primitives_test

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"testing"

	"github.com/hyperledger/fabric/core/crypto/primitives"
)

func TestCBCPKCS7EncryptCBCPKCS7Decrypt(t *testing.T) {
	// Encrypt with CBCPKCS7Encrypt and Decrypt with CBCPKCS7Decrypt
	// The purpose of this test is not to test the implementation of the AES standard
	// library but to verify the code wrapping/unwrapping the AES cipher.

	key := make([]byte, primitives.AESKeyLength)
	rand.Reader.Read(key)

	var msg = []byte("a message with arbitrary length (42 bytes)")

	encrypted, encErr := primitives.CBCPKCS7Encrypt(key, msg)

	if encErr != nil {
		t.Fatalf("Error encrypting message %v", encErr)
	}

	decrypted, dErr := primitives.CBCPKCS7Decrypt(key, encrypted)

	if dErr != nil {
		t.Fatalf("Error encrypting message %v", dErr)
	}

	if string(msg[:]) != string(decrypted[:]) {
		t.Fatalf("Decrypt( Encrypt(msg) ) != msg: Ciphertext decryption with the same key must restore the original message!")
	}

}

func TestPKCS7Padding(t *testing.T) {
	// Verify the PKCS7 padding, using a human readable plaintext.

	// 0 byte/length message
	msg := []byte("")
	expected := []byte{16, 16, 16, 16,
		16, 16, 16, 16,
		16, 16, 16, 16,
		16, 16, 16, 16}
	result := primitives.PKCS7Padding(msg)

	if !bytes.Equal(expected, result) {
		t.Fatal("Padding error! Expected: ", expected, "', received: '", result, "'")
	}

	// 1 byte/length message
	msg = []byte("0")
	expected = []byte{'0', 15, 15, 15,
		15, 15, 15, 15,
		15, 15, 15, 15,
		15, 15, 15, 15}
	result = primitives.PKCS7Padding(msg)

	if !bytes.Equal(expected, result) {
		t.Fatal("Padding error! Expected: '", expected, "', received: '", result, "'")
	}

	// 2 byte/length message
	msg = []byte("01")
	expected = []byte{'0', '1', 14, 14,
		14, 14, 14, 14,
		14, 14, 14, 14,
		14, 14, 14, 14}
	result = primitives.PKCS7Padding(msg)

	if !bytes.Equal(expected, result) {
		t.Fatal("Padding error! Expected: '", expected, "', received: '", result, "'")
	}

	// 3 to aes.BlockSize-1 byte messages
	for i := 3; i < aes.BlockSize; i++ {
		msg := []byte("0123456789ABCDEF")

		result := primitives.PKCS7Padding(msg[:i])

		padding := aes.BlockSize - i
		expectedPadding := bytes.Repeat([]byte{byte(padding)}, padding)
		expected = append(msg[:i], expectedPadding...)

		if !bytes.Equal(result, expected) {
			t.Fatal("Padding error! Expected: '", expected, "', received: '", result, "'")
		}

	}

	// aes.BlockSize length message
	msg = bytes.Repeat([]byte{byte('x')}, aes.BlockSize)

	result = primitives.PKCS7Padding(msg)

	expectedPadding := bytes.Repeat([]byte{byte(aes.BlockSize)},
		aes.BlockSize)
	expected = append(msg, expectedPadding...)

	if len(result) != 2*aes.BlockSize {
		t.Fatal("Padding error: expected the length of the returned slice to be 2 times aes.BlockSize")
	}

	if !bytes.Equal(expected, result) {
		t.Fatal("Padding error! Expected: '", expected, "', received: '", result, "'")
	}

}

func TestPKCS7UnPadding(t *testing.T) {
	// Verify the PKCS7 unpadding, using a human readable plaintext.

	// 0 byte/length message
	expected := []byte("")
	msg := []byte{16, 16, 16, 16,
		16, 16, 16, 16,
		16, 16, 16, 16,
		16, 16, 16, 16}

	result, _ := primitives.PKCS7UnPadding(msg)

	if !bytes.Equal(expected, result) {
		t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
	}

	// 1 byte/length message
	expected = []byte("0")
	msg = []byte{'0', 15, 15, 15,
		15, 15, 15, 15,
		15, 15, 15, 15,
		15, 15, 15, 15}

	result, _ = primitives.PKCS7UnPadding(msg)

	if !bytes.Equal(expected, result) {
		t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
	}

	// 2 byte/length message
	expected = []byte("01")
	msg = []byte{'0', '1', 14, 14,
		14, 14, 14, 14,
		14, 14, 14, 14,
		14, 14, 14, 14}

	result, _ = primitives.PKCS7UnPadding(msg)

	if !bytes.Equal(expected, result) {
		t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
	}

	// 3 to aes.BlockSize-1 byte messages
	for i := 3; i < aes.BlockSize; i++ {
		base := []byte("0123456789ABCDEF")

		iPad := aes.BlockSize - i
		padding := bytes.Repeat([]byte{byte(iPad)}, iPad)
		msg = append(base[:i], padding...)

		expected := base[:i]
		result, _ := primitives.PKCS7UnPadding(msg)

		if !bytes.Equal(result, expected) {
			t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
		}

	}

	// aes.BlockSize length message
	expected = bytes.Repeat([]byte{byte('x')}, aes.BlockSize)

	padding := bytes.Repeat([]byte{byte(aes.BlockSize)},
		aes.BlockSize)
	msg = append(expected, padding...)

	result, _ = primitives.PKCS7UnPadding(msg)

	if !bytes.Equal(expected, result) {
		t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
	}
}

func TestCBCEncryptCBCPKCS7Decrypt_ExpectingFailure(t *testing.T) {
	// When encrypting a message that does not need padding, i.e. a message
	// whose length is a multiple of the block size, with CBCEncrypt, it
	// cannot be decrypted with CBCPKCS7DEcrypt.
	//
	// The intend of this test is to document this behaviour for clarity.
	//
	// The reason is that the section 10.3 Note 2 in #PKCS7 states that
	// an extra block is appended to the message for padding. Since this
	// extra block is missing when using CBCEncrypt for encryption,
	// CBCPKCS7Decrypt fails.

	key := make([]byte, 32)
	rand.Reader.Read(key)

	var msg = []byte("a 16 byte messag")

	encrypted, encErr := CBCEncrypt(key, msg)

	if encErr != nil {
		t.Fatalf("Error encrypting message %v", encErr)
	}

	decrypted, dErr := CBCPKCS7Decrypt(key, encrypted)

	if dErr != nil {
		// expected behaviour, decryption fails.
		t.Logf("Expected error decrypting message %v, decrypted message = %v", dErr, decrypted)
	}
}

func TestCBCPKCS7EncryptCBCDecrypt_ExpectingCorruptMessage(t *testing.T) {
	// When encrypting a message that does not need padding, i.e. a message
	// whose length is a multiple of the block size, with CBCPKCS7Encrypt, it
	// can be decrypted with CBCDecrypt but the returned message is corrupted.
	//
	// The intend of this test is to document this behaviour for clarity.
	//
	// The reason is that the section 10.3 Note 2 in #PKCS7 states that
	// an extra block is appended to the message for padding. Since this
	// extra block is added when using CBCPKCS7Encrypt for encryption,
	// CBCDecrypt returns the original message plus this extra block.
	//
	// Note:
	// the same applies for messages of arbitrary length.

	key := make([]byte, 32)
	rand.Reader.Read(key)

	//                0123456789ABCDEF
	var msg = []byte("a 16 byte messag")

	encrypted, encErr := CBCPKCS7Encrypt(key, msg)

	if encErr != nil {
		t.Fatalf("Error encrypting message %v", encErr)
	}

	decrypted, dErr := CBCDecrypt(key, encrypted)

	if dErr != nil {
		t.Fatalf("Error encrypting message %v, %v", dErr, decrypted)
	}

	if string(msg[:]) != string(decrypted[:aes.BlockSize]) {
		t.Log("msg: ", msg)
		t.Log("decrypted: ", decrypted[:aes.BlockSize])
		t.Fatalf("Encryption->Decryption with same key should result in original message")
	}

	if !bytes.Equal(decrypted[aes.BlockSize:], bytes.Repeat([]byte{byte(aes.BlockSize)}, aes.BlockSize)) {
		t.Fatal("Expected extra block with padding in encrypted message", decrypted)
	}

}

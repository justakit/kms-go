/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package subtle_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"testing"

	josecipher "github.com/go-jose/go-jose/v3/cipher"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/kms-crypto-go/crypto/tinkcrypto/primitive/aead/subtle"
)

func TestNewAESCBCHMAC(t *testing.T) {
	key := make([]byte, 64)

	// Test various key sizes.
	for i := 0; i < 64; i++ {
		k := key[:i]
		keySize := len(k)

		c, err := subtle.NewAESCBCHMAC(k)

		switch keySize {
		case 32, 48, 64:
			// Valid key sizes.
			require.NoError(t, err, "want: valid cipher (key size=%d), got: error %v", len(k), err)

			// Verify that the struct contents are correctly set.
			require.Equal(t, len(k), len(c.Key), "want: key size=%d, got: key size=%d", keySize, len(c.Key))
		default:
			require.EqualError(t, err, fmt.Sprintf("aes_cbc_hmac: invalid AES CBC key size; want 32, 48 or 64, got %d", keySize))
		}
	}
}

func TestIETFTestVector(t *testing.T) {
	// Source: https://tools.ietf.org/html/draft-mcgrew-aead-aes-cbc-hmac-sha2-05#section-5
	plaintext := []byte{
		0x41, 0x20, 0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x20, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x20,
		0x6d, 0x75, 0x73, 0x74, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x62, 0x65, 0x20, 0x72, 0x65, 0x71, 0x75,
		0x69, 0x72, 0x65, 0x64, 0x20, 0x74, 0x6f, 0x20, 0x62, 0x65, 0x20, 0x73, 0x65, 0x63, 0x72, 0x65,
		0x74, 0x2c, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x69, 0x74, 0x20, 0x6d, 0x75, 0x73, 0x74, 0x20, 0x62,
		0x65, 0x20, 0x61, 0x62, 0x6c, 0x65, 0x20, 0x74, 0x6f, 0x20, 0x66, 0x61, 0x6c, 0x6c, 0x20, 0x69,
		0x6e, 0x74, 0x6f, 0x20, 0x74, 0x68, 0x65, 0x20, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x20, 0x6f, 0x66,
		0x20, 0x74, 0x68, 0x65, 0x20, 0x65, 0x6e, 0x65, 0x6d, 0x79, 0x20, 0x77, 0x69, 0x74, 0x68, 0x6f,
		0x75, 0x74, 0x20, 0x69, 0x6e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x6e, 0x69, 0x65, 0x6e, 0x63, 0x65,
	}

	aad := []byte{
		0x54, 0x68, 0x65, 0x20, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x20, 0x70, 0x72, 0x69, 0x6e, 0x63,
		0x69, 0x70, 0x6c, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x41, 0x75, 0x67, 0x75, 0x73, 0x74, 0x65, 0x20,
		0x4b, 0x65, 0x72, 0x63, 0x6b, 0x68, 0x6f, 0x66, 0x66, 0x73,
	}

	nonce := []byte{
		0x1a, 0xf3, 0x8c, 0x2d, 0xc2, 0xb9, 0x6f, 0xfd, 0xd8, 0x66, 0x94, 0x09, 0x23, 0x41, 0xbc, 0x04,
	}

	expectedCiphertext1 := []byte{
		0xc8, 0x0e, 0xdf, 0xa3, 0x2d, 0xdf, 0x39, 0xd5, 0xef, 0x00, 0xc0, 0xb4, 0x68, 0x83, 0x42, 0x79,
		0xa2, 0xe4, 0x6a, 0x1b, 0x80, 0x49, 0xf7, 0x92, 0xf7, 0x6b, 0xfe, 0x54, 0xb9, 0x03, 0xa9, 0xc9,
		0xa9, 0x4a, 0xc9, 0xb4, 0x7a, 0xd2, 0x65, 0x5c, 0x5f, 0x10, 0xf9, 0xae, 0xf7, 0x14, 0x27, 0xe2,
		0xfc, 0x6f, 0x9b, 0x3f, 0x39, 0x9a, 0x22, 0x14, 0x89, 0xf1, 0x63, 0x62, 0xc7, 0x03, 0x23, 0x36,
		0x09, 0xd4, 0x5a, 0xc6, 0x98, 0x64, 0xe3, 0x32, 0x1c, 0xf8, 0x29, 0x35, 0xac, 0x40, 0x96, 0xc8,
		0x6e, 0x13, 0x33, 0x14, 0xc5, 0x40, 0x19, 0xe8, 0xca, 0x79, 0x80, 0xdf, 0xa4, 0xb9, 0xcf, 0x1b,
		0x38, 0x4c, 0x48, 0x6f, 0x3a, 0x54, 0xc5, 0x10, 0x78, 0x15, 0x8e, 0xe5, 0xd7, 0x9d, 0xe5, 0x9f,
		0xbd, 0x34, 0xd8, 0x48, 0xb3, 0xd6, 0x95, 0x50, 0xa6, 0x76, 0x46, 0x34, 0x44, 0x27, 0xad, 0xe5,
		0x4b, 0x88, 0x51, 0xff, 0xb5, 0x98, 0xf7, 0xf8, 0x00, 0x74, 0xb9, 0x47, 0x3c, 0x82, 0xe2, 0xdb,
	}

	expectedAuthtag1 := []byte{
		0x65, 0x2c, 0x3f, 0xa3, 0x6b, 0x0a, 0x7c, 0x5b, 0x32, 0x19, 0xfa, 0xb3, 0xa3, 0x0b, 0xc1, 0xc4,
	}

	key1 := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}

	expectedCiphertext2 := []byte{
		0xea, 0x65, 0xda, 0x6b, 0x59, 0xe6, 0x1e, 0xdb, 0x41, 0x9b, 0xe6, 0x2d, 0x19, 0x71, 0x2a, 0xe5,
		0xd3, 0x03, 0xee, 0xb5, 0x00, 0x52, 0xd0, 0xdf, 0xd6, 0x69, 0x7f, 0x77, 0x22, 0x4c, 0x8e, 0xdb,
		0x00, 0x0d, 0x27, 0x9b, 0xdc, 0x14, 0xc1, 0x07, 0x26, 0x54, 0xbd, 0x30, 0x94, 0x42, 0x30, 0xc6,
		0x57, 0xbe, 0xd4, 0xca, 0x0c, 0x9f, 0x4a, 0x84, 0x66, 0xf2, 0x2b, 0x22, 0x6d, 0x17, 0x46, 0x21,
		0x4b, 0xf8, 0xcf, 0xc2, 0x40, 0x0a, 0xdd, 0x9f, 0x51, 0x26, 0xe4, 0x79, 0x66, 0x3f, 0xc9, 0x0b,
		0x3b, 0xed, 0x78, 0x7a, 0x2f, 0x0f, 0xfc, 0xbf, 0x39, 0x04, 0xbe, 0x2a, 0x64, 0x1d, 0x5c, 0x21,
		0x05, 0xbf, 0xe5, 0x91, 0xba, 0xe2, 0x3b, 0x1d, 0x74, 0x49, 0xe5, 0x32, 0xee, 0xf6, 0x0a, 0x9a,
		0xc8, 0xbb, 0x6c, 0x6b, 0x01, 0xd3, 0x5d, 0x49, 0x78, 0x7b, 0xcd, 0x57, 0xef, 0x48, 0x49, 0x27,
		0xf2, 0x80, 0xad, 0xc9, 0x1a, 0xc0, 0xc4, 0xe7, 0x9c, 0x7b, 0x11, 0xef, 0xc6, 0x00, 0x54, 0xe3,
	}

	expectedAuthtag2 := []byte{
		0x84, 0x90, 0xac, 0x0e, 0x58, 0x94, 0x9b, 0xfe, 0x51, 0x87, 0x5d, 0x73, 0x3f, 0x93, 0xac, 0x20,
		0x75, 0x16, 0x80, 0x39, 0xcc, 0xc7, 0x33, 0xd7,
	}

	key2 := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
	}

	// Key3 is not 32, 48 or 64 in size and therefore not supported by go-jose.
	// key3 := []byte{
	//	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	//	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	//	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
	//	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37}

	expectedCiphertext4 := []byte{
		0x4a, 0xff, 0xaa, 0xad, 0xb7, 0x8c, 0x31, 0xc5, 0xda, 0x4b, 0x1b, 0x59, 0x0d, 0x10, 0xff, 0xbd,
		0x3d, 0xd8, 0xd5, 0xd3, 0x02, 0x42, 0x35, 0x26, 0x91, 0x2d, 0xa0, 0x37, 0xec, 0xbc, 0xc7, 0xbd,
		0x82, 0x2c, 0x30, 0x1d, 0xd6, 0x7c, 0x37, 0x3b, 0xcc, 0xb5, 0x84, 0xad, 0x3e, 0x92, 0x79, 0xc2,
		0xe6, 0xd1, 0x2a, 0x13, 0x74, 0xb7, 0x7f, 0x07, 0x75, 0x53, 0xdf, 0x82, 0x94, 0x10, 0x44, 0x6b,
		0x36, 0xeb, 0xd9, 0x70, 0x66, 0x29, 0x6a, 0xe6, 0x42, 0x7e, 0xa7, 0x5c, 0x2e, 0x08, 0x46, 0xa1,
		0x1a, 0x09, 0xcc, 0xf5, 0x37, 0x0d, 0xc8, 0x0b, 0xfe, 0xcb, 0xad, 0x28, 0xc7, 0x3f, 0x09, 0xb3,
		0xa3, 0xb7, 0x5e, 0x66, 0x2a, 0x25, 0x94, 0x41, 0x0a, 0xe4, 0x96, 0xb2, 0xe2, 0xe6, 0x60, 0x9e,
		0x31, 0xe6, 0xe0, 0x2c, 0xc8, 0x37, 0xf0, 0x53, 0xd2, 0x1f, 0x37, 0xff, 0x4f, 0x51, 0x95, 0x0b,
		0xbe, 0x26, 0x38, 0xd0, 0x9d, 0xd7, 0xa4, 0x93, 0x09, 0x30, 0x80, 0x6d, 0x07, 0x03, 0xb1, 0xf6,
	}

	expectedAuthtag4 := []byte{
		0x4d, 0xd3, 0xb4, 0xc0, 0x88, 0xa7, 0xf4, 0x5c, 0x21, 0x68, 0x39, 0x64, 0x5b, 0x20, 0x12, 0xbf,
		0x2e, 0x62, 0x69, 0xa8, 0xc5, 0x6a, 0x81, 0x6d, 0xbc, 0x1b, 0x26, 0x77, 0x61, 0x95, 0x5b, 0xc5,
	}

	key4 := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
	}

	tests := []struct {
		name               string
		plaintext          []byte
		aad                []byte
		expectedCiphertext []byte
		expectedAuthtag    []byte
		key                []byte
		nonce              []byte
	}{
		{
			name:               "AEAD_AES_128_CBC_HMAC_SHA256",
			plaintext:          plaintext,
			aad:                aad,
			expectedCiphertext: expectedCiphertext1,
			expectedAuthtag:    expectedAuthtag1,
			key:                key1,
			nonce:              nonce,
		},
		{
			name:               "AEAD_AES_192_CBC_HMAC_SHA384",
			plaintext:          plaintext,
			aad:                aad,
			expectedCiphertext: expectedCiphertext2,
			expectedAuthtag:    expectedAuthtag2,
			key:                key2,
			nonce:              nonce,
		},
		// {
		//	name:               "AEAD_AES_256_CBC_HMAC_SHA384",
		//	plaintext:          plaintext,
		//	aad:                aad,
		//	expectedCiphertext: expectedCiphertext3,
		//	expectedAuthtag:    expectedAuthtag3,
		// Key3 is not supported by Go-Jose (key length=56 not supported). This is why this test is commented out.
		//	key:                key3,
		//	nonce:              nonce,
		// },
		{
			name:               "AEAD_AES_256_CBC_HMAC_SHA512",
			plaintext:          plaintext,
			aad:                aad,
			expectedCiphertext: expectedCiphertext4,
			expectedAuthtag:    expectedAuthtag4,
			key:                key4,
			nonce:              nonce,
		},
	}

	t.Parallel()

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			cbcHMAC, err := josecipher.NewCBCHMAC(tc.key, aes.NewCipher)
			require.NoError(t, err)

			enc := mockNONCEInCBCHMAC{
				nonce:   nonce,
				cbcHMAC: cbcHMAC,
			}

			out, err := enc.Encrypt(plaintext, aad)
			require.NoError(t, err, "unable to encrypt")

			tagSize := len(tc.expectedAuthtag)

			ct := make([]byte, len(nonce)+len(tc.expectedCiphertext)+len(tc.expectedAuthtag))
			copy(ct, nonce)
			copy(ct[len(nonce):], tc.expectedCiphertext)
			copy(ct[len(nonce)+len(tc.expectedCiphertext):], tc.expectedAuthtag)

			out1, err := enc.Decrypt(ct, aad)
			require.NoError(t, err, "unable to decrypt")

			require.EqualValues(t, plaintext, out1)

			if !bytes.Equal(out[len(nonce):len(out)-tagSize], tc.expectedCiphertext) {
				t.Error("Ciphertext did not match, got", out[len(nonce):len(out)-tagSize], "wanted", tc.expectedCiphertext)
			}

			if !bytes.Equal(out[len(out)-tagSize:], tc.expectedAuthtag) {
				t.Error("Auth tag did not match, got", out[len(out)-tagSize:], "wanted", tc.expectedAuthtag)
			}
		})
	}
}

type mockNONCEInCBCHMAC struct {
	subtle.AESCBCHMAC

	cbcHMAC cipher.AEAD
	nonce   []byte
}

// Encrypt using the mocked nonce instead of generating a random one.
func (a *mockNONCEInCBCHMAC) Encrypt(plaintext, additionalData []byte) ([]byte, error) {
	AESCBCIVSize := 16

	ciphertext := a.cbcHMAC.Seal(nil, a.nonce, plaintext, additionalData)

	ciphertextAndIV := make([]byte, AESCBCIVSize+len(ciphertext))
	if n := copy(ciphertextAndIV, a.nonce); n != AESCBCIVSize {
		return nil, fmt.Errorf("aes_cbc_hmac: failed to copy IV (copied %d/%d bytes)", n, AESCBCIVSize)
	}

	copy(ciphertextAndIV[AESCBCIVSize:], ciphertext)

	return ciphertextAndIV, nil
}

func (a *mockNONCEInCBCHMAC) Decrypt(ciphertext, additionalData []byte) ([]byte, error) {
	ivSize := a.cbcHMAC.NonceSize()
	if len(ciphertext) < ivSize {
		return nil, fmt.Errorf("aes_cbc_hmac: ciphertext too short")
	}

	iv := ciphertext[:ivSize]

	return a.cbcHMAC.Open(nil, iv, ciphertext[ivSize:], additionalData)
}

func TestAESCBCRoundtrip(t *testing.T) {
	key128 := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	}

	key192 := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
	}

	key256 := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	}

	RunRoundtrip(t, key128)
	RunRoundtrip(t, key192)
	RunRoundtrip(t, key256)
}

func RunRoundtrip(t *testing.T, key []byte) {
	aead, err := subtle.NewAESCBCHMAC(key)
	require.NoError(t, err)

	// Test pre-existing data in dst buffer
	plaintext := []byte{0, 0, 0, 0}
	aad := []byte{4, 3, 2, 1}

	result, err := aead.Encrypt(plaintext, aad)
	require.NoError(t, err)

	result, err = aead.Decrypt(result, aad)
	require.NoError(t, err)
	require.EqualValues(t, plaintext, result, "Plaintext does not match output")

	t.Run("failure: bad cipher", func(t *testing.T) {
		result, err = aead.Decrypt([]byte("bad cipher"), aad)
		require.EqualError(t, err, "aes_cbc_hmac: ciphertext too short")
	})

	t.Run("failure: cipher not short but not large enough to contain an authentication tag", func(t *testing.T) {
		result, err = aead.Decrypt([]byte("bad cipher with not too short length to cause decryption failure"), aad)
		require.EqualError(t, err, "aes_cbc_hmac: failed to decrypt: go-jose/go-jose: invalid ciphertext "+
			"(auth tag mismatch)")
	})
}

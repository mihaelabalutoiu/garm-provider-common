package util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	EncryptionPassphrase     = "bocyasicgatEtenOubwonIbsudNutDom"
	WeakEncryptionPassphrase = "test123"
)

func TestSealUnseal(t *testing.T) {
	data, err := Seal([]byte("test"), []byte(EncryptionPassphrase))
	require.NoError(t, err)

	// Data should unmarshal in an Envelope
	var envelope Envelope
	err = json.Unmarshal(data, &envelope)
	require.NoError(t, err)

	// Should be able to decrypt
	decrypted, err := Unseal(data, []byte(EncryptionPassphrase))
	require.NoError(t, err)
	require.Equal(t, []byte("test"), decrypted)
}

func TestAes256EncodeDecode(t *testing.T) {
	// Should be able to encrypt
	encrypted, err := Aes256Encode([]byte("test"), EncryptionPassphrase)
	require.NoError(t, err)

	// Should be able to decrypt
	decrypted, err := Aes256Decode(encrypted, EncryptionPassphrase)
	require.NoError(t, err)
	require.Equal(t, []byte("test"), decrypted)
}

func TestAes256EncodeDecodeWeakSecret(t *testing.T) {
	_, err := Aes256Encode([]byte("test"), WeakEncryptionPassphrase)
	require.NotNil(t, err)
	require.EqualError(t, err, "invalid passphrase length (expected length 32 characters)")

	_, err = Aes256Decode([]byte("test"), WeakEncryptionPassphrase)
	require.NotNil(t, err)
	require.EqualError(t, err, "invalid passphrase length (expected length 32 characters)")
}

func TestUnsealAes256EncodedData(t *testing.T) {
	encrypted, err := Aes256Encode([]byte("test"), EncryptionPassphrase)
	require.NoError(t, err)

	// Should be able to decrypt
	decrypted, err := Unseal(encrypted, []byte(EncryptionPassphrase))
	require.NoError(t, err)
	require.Equal(t, []byte("test"), decrypted)
}

func TestSealUnsealWeakSecret(t *testing.T) {
	_, err := Seal([]byte("test"), []byte(WeakEncryptionPassphrase))
	require.NotNil(t, err)
	require.EqualError(t, err, "invalid passphrase length (expected length 32 characters)")

	// The data is irrelevant. We expect to error out on the passphrase length.
	_, err = Unseal([]byte("test"), []byte(WeakEncryptionPassphrase))
	require.NotNil(t, err)
	require.EqualError(t, err, "invalid passphrase length (expected length 32 characters)")
}

func TestAes256EncodeDecodeString(t *testing.T) {
	encrypted, err := Aes256EncodeString("test", EncryptionPassphrase)
	require.NoError(t, err)

	decrypted, err := Aes256DecodeString(encrypted, EncryptionPassphrase)
	require.NoError(t, err)
	require.Equal(t, "test", decrypted)
}

func TestAes256EncodeStringWeakSecret(t *testing.T) {
	_, err := Aes256EncodeString("test", WeakEncryptionPassphrase)
	require.NotNil(t, err)
	require.EqualError(t, err, "invalid passphrase length (expected length 32 characters)")
}

func TestAes256DecodeWrongEncryptedString(t *testing.T) {
	_, err := Aes256DecodeString([]byte(""), EncryptionPassphrase)
	require.NotNil(t, err)
	require.EqualError(t, err, "failed to decrypt text")
}

func TestAes256DecodeWrongDecryptionPassphrase(t *testing.T) {
	encrypted, err := Aes256EncodeString("test", EncryptionPassphrase)
	require.NoError(t, err)

	// We pass a wrong decryption passphrase, that it's still 32 characters long.
	_, err = Aes256DecodeString(encrypted, "wrong passphrase-1234-1234-12345")
	require.NotNil(t, err)
	require.EqualError(t, err, "failed to decrypt text")
}

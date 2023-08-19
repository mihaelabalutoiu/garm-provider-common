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

	// The data is irelevant. We expect to error out on the passphrase length.
	_, err = Unseal([]byte("test"), []byte(WeakEncryptionPassphrase))
	require.NotNil(t, err)
	require.EqualError(t, err, "invalid passphrase length (expected length 32 characters)")
}

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/OGLinuk/go-sbh/sbh"
)

func genSBH(passphrase string, nrots, seed int64) string {
	return sbh.Generate(passphrase, nrots, seed)
}

func getBlockCipher(SBH string) (cipher.AEAD, error) {
	block, err := aes.NewCipher([]byte(SBH))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}

func encrypt(data []byte, SBH string) ([]byte, error) {
	gcm, err := getBlockCipher(SBH[:32])
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted := gcm.Seal(nonce, nonce, data, nil)

	return encrypted, nil
}

func encryptFile(data []byte, w io.Writer, SBH string) error {
	bytes, err := encrypt(data, SBH)
	if err != nil {
		return err
	}

	w.Write(bytes)

	return nil
}

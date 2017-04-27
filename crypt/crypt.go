package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
)

type encodedMessage struct {
	Key     []byte
	Nonce   []byte
	Content []byte
}

func GenerateKey(strength int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, strength)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func Encode(msg string, key *rsa.PublicKey) (string, error) {

	//First generate a random AES key
	secret := make([]byte, 32)
	_, _ = rand.Read(secret)
	enc := &encodedMessage{}

	//And encrypt the message with the AES key
	nonce := make([]byte, 12)
	rand.Read(nonce)
	block, _ := aes.NewCipher(secret)
	gcm, _ := cipher.NewGCM(block)
	enc.Content = gcm.Seal(nil, nonce, []byte(msg), nil)

	var err error
	//Then encrypt the key and nonce with their public key
	enc.Key, _ = rsa.EncryptOAEP(sha256.New(), rand.Reader, key, secret, []byte("key"))
	enc.Nonce, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, key, nonce, []byte("nonce"))

	//Finally encode the struct into a string
	res := &bytes.Buffer{}
	gobber := gob.NewEncoder(res)
	gobber.Encode(enc)

	final := res.String()

	return final, err
}

func Decode(msg string, key *rsa.PrivateKey) (string, error) {

	//First decode the struct
	reader := bytes.NewBufferString(msg)
	gobber := gob.NewDecoder(reader)
	enc := &encodedMessage{}
	gobber.Decode(enc)

	//Then decrypt the AES key with our private key
	aes_key, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, enc.Key, []byte("key"))
	nonce, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, enc.Nonce, []byte("nonce"))

	//Finally decrypt the Content with the AES key
	block, _ := aes.NewCipher(aes_key)
	gcm, _ := cipher.NewGCM(block)
	res, _ := gcm.Open(nil, nonce, enc.Content, nil)

	//And return it as a string
	return string(res), err
}

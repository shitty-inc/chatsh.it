package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"syscall/js"

	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
	"github.com/keybase/go-triplesec"
	"github.com/shitty-inc/sendshit-go"
	"golang.org/x/crypto/curve25519"
)

var PrivateKey [32]byte

func GenerateRandomString(this js.Value, args []js.Value) (interface{}, error) {
	return sendshit.GenerateRandomString(args[0].Int())
}

func Encrypt(this js.Value, args []js.Value) (interface{}, error) {
	cipher, err := triplesec.NewCipher([]byte(args[2].String()), nil, 3)

	if err != nil {
		return nil, err
	}

	size := args[1].Int()
	bytes := make([]byte, size)

	js.CopyBytesToGo(bytes, args[0])

	encrypted, err := cipher.Encrypt(bytes)
	if err != nil {
		return nil, err
	}

	data := js.Global().Get("Uint8Array").New(len(encrypted))
	js.CopyBytesToJS(data, encrypted)

	return js.ValueOf(data), nil
}

func Decrypt(this js.Value, args []js.Value) (interface{}, error) {
	cipher, err := triplesec.NewCipher([]byte(args[2].String()), nil, 3)

	if err != nil {
		return nil, err
	}

	size := args[1].Int()
	bytes := make([]byte, size)

	js.CopyBytesToGo(bytes, args[0])

	decrypted, err := cipher.Decrypt(bytes)
	if err != nil {
		return nil, err
	}

	data := js.Global().Get("Uint8Array").New(len(decrypted))
	js.CopyBytesToJS(data, decrypted)

	return js.ValueOf(data), nil
}

func GenerateKey(this js.Value, args []js.Value) (interface{}, error) {
	var privKey [32]byte

	_, err := io.ReadFull(rand.Reader, privKey[:])

	privKey[0] &= 248
	privKey[31] &= 127
	privKey[31] |= 64

	var pubKey [32]byte

	curve25519.ScalarBaseMult(&pubKey, &privKey)

	PrivateKey = privKey

	return js.ValueOf(hex.EncodeToString(pubKey[:])), err
}

func ComputeSecret(this js.Value, args []js.Value) (interface{}, error) {

	publicString := args[0].String()
	publicBytes, _ := hex.DecodeString(publicString)

	var secret [32]byte

	curve25519.ScalarMult(&secret, &PrivateKey, (*[32]byte)(publicBytes) )

	return js.ValueOf(hex.EncodeToString(secret[:])), nil
}

func main() {
	c := make(chan struct{}, 0)

	gobridge.RegisterCallback("GenerateRandomString", GenerateRandomString)
	gobridge.RegisterCallback("Encrypt", Encrypt)
	gobridge.RegisterCallback("Decrypt", Decrypt)
	gobridge.RegisterCallback("GenerateKey", GenerateKey)
	gobridge.RegisterCallback("ComputeSecret", ComputeSecret)

	<-c
}

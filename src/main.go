package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"syscall/js"

	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
	"github.com/shitty-inc/sendshit-go"
	"golang.org/x/crypto/curve25519"
)

var PrivateKey [32]byte

func GenerateRandomString(this js.Value, args []js.Value) (interface{}, error) {
	return sendshit.GenerateRandomString(args[0].Int())
}

func EncryptFile(this js.Value, args []js.Value) (interface{}, error) {
	size := args[2].Int()
	image := make([]byte, size)
	js.CopyBytesToGo(image, args[1])

	return sendshit.EncryptFile(args[0].String(), image, args[3].String())
}

func UploadFile(this js.Value, args []js.Value) (interface{}, error) {
	value, err := sendshit.UploadFile(args[0].String())

	return js.ValueOf(value), err
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

	return js.ValueOf(base64.StdEncoding.EncodeToString(pubKey[:])), err
}

func ComputeSecret(this js.Value, args []js.Value) (interface{}, error) {

	publicString := args[0].String()
	publicBytes, _ := base64.StdEncoding.DecodeString(publicString)

	var secret [32]byte

	curve25519.ScalarMult(&secret, &PrivateKey, (*[32]byte)(publicBytes) )

	return js.ValueOf(base64.StdEncoding.EncodeToString(secret[:])), nil
}

func main() {
	c := make(chan struct{}, 0)

	gobridge.RegisterCallback("GenerateRandomString", GenerateRandomString)
	gobridge.RegisterCallback("EncryptFile", EncryptFile)
	gobridge.RegisterCallback("UploadFile", UploadFile)
	gobridge.RegisterCallback("GenerateKey", GenerateKey)
	gobridge.RegisterCallback("ComputeSecret", ComputeSecret)

	<-c
}

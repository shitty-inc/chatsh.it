package main

import (
	"syscall/js"

	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
	"github.com/shitty-inc/sendshit-go"
)

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

func main() {
	c := make(chan struct{}, 0)

	gobridge.RegisterCallback("GenerateRandomString", GenerateRandomString)
	gobridge.RegisterCallback("EncryptFile", EncryptFile)
	gobridge.RegisterCallback("UploadFile", UploadFile)

	<-c
}

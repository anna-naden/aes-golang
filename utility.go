package main

import (
	"encoding/base64"
	"io/ioutil"
)

func get_cipher_bytes(path string) []byte {
	base64sftr, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	b64 := string(base64sftr)
	cipher_text_bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	return cipher_text_bytes
}

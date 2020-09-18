package main

import (
	// "encoding/hex"
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

type STATE [4][4]uint8

func main() {
	const key = "YELLOW SUBMARINE"
	keys1 := []byte(key)
	keys := [16]uint8{}
	for i := 0; i < 16; i++ {
		keys[i] = keys1[i]
	}
	base64str, err := ioutil.ReadFile("crypto-challenge.txt")
	if err != nil {
		panic(err)
	}
	b64 := string(base64str)
	bytes1, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	p2 := [16]uint8{}
	fmt.Println(len(bytes1))
	f, err := os.Create("plaintext.txt")
	w := bufio.NewWriter(f)
	n := 0
	for run := 0; run < 1; run++ {
		for j := 0; j+15 < len(bytes1); j += 16 {
			for i := 0; i < 16; i++ {
				p2[i] = uint8(bytes1[i+j])
			}
			state := decrypt(initializeState(p2), key_expand(keys))
			bytes := unpackState(state)
			n += len(bytes)
			w.WriteString(string(bytes[:]))
			fmt.Println(string(bytes[:]))
		}
	}
	w.Flush()
	fmt.Println(n)
}

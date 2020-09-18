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

type SBOX_LOOKER_UPPER interface{
	lookup() STATE
}

type PLAIN_TEXT_STATE STATE

type CIPHER_TEXT_STATE STATE

func main() {

	const key = "YELLOW SUBMARINE"
	keys1 := []byte(key)
	key_bytes := [16]uint8{}
	for i := 0; i < 16; i++ {
		key_bytes[i] = keys1[i]
	}
	base64str, err := ioutil.ReadFile("crypto-challenge.txt")
	if err != nil {
		panic(err)
	}
	b64 := string(base64str)
	cipher_text_bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	block := [16]uint8{} //Blocks of ciphertext
	fmt.Println(len(cipher_text_bytes))
	f, err := os.Create("plaintext.txt")
	w := bufio.NewWriter(f)
	n := 0
	for run := 0; run < 1; run++ {
		for j := 0; j+15 < len(cipher_text_bytes); j += 16 {
			for i := 0; i < 16; i++ {
				block[i] = uint8(cipher_text_bytes[i+j])
			}
			state := decrypt(initializeState(block), get_key_schedule(key_bytes))
			plain_text_bytes := unpackState(state)
			n += len(plain_text_bytes)
			w.WriteString(string(plain_text_bytes[:]))
			fmt.Println(string(plain_text_bytes[:]))
		}
	}
	w.Flush()
	fmt.Println(n)
}

func decrypt(cipher_matrix STATE, key_schedule [44]uint32) STATE {

	//Round 0
	round := 0
	round_key := make_key_matrix(key_schedule[40-round*4 : 44-round*4])
	cipher_matrix.add_round_key(round_key)
	state := cipher_matrix

	for round := 1; round <= 9; round++ {
		state.inv_shift_rows()
		state.substitute_inv_sbox()
		key_matrix := make_key_matrix(key_schedule[40-4*round : 44-4*round])
		state.add_round_key(key_matrix)
		state.inv_mix_columns()
	}

	round = 10
	state.inv_shift_rows()
	state.substitute_inv_sbox()
	round_key = make_key_matrix(key_schedule[40-4*round : 44-4*round])
	state.add_round_key(round_key)
	return state
}
func encrypt(plaintext [16]uint8, key_schedule [44]uint32) STATE {

	round_key := make_key_matrix(key_schedule[0:4])
	state := initializeState(plaintext)
	state.add_round_key(round_key)

	for round := 1; round <= 9; round++ {
		state.substitute_sbox()
		state.shift_rows()
		state.MixColumns()
		key_matrix := make_key_matrix(key_schedule[4*round : 4*round+4])
		state.add_round_key(key_matrix)
	}

	round := 10
	state.substitute_sbox()
	state.shift_rows()
	round_key = make_key_matrix(key_schedule[4*round : 4*round+4])
	state.add_round_key(round_key)
	return state
}

package main

import (
	// "encoding/hex"
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

type STATE [4][4]byte

type SBOX_LOOKER_UPPER interface {
	lookup() STATE
}

type ENCRYPTION_STATE struct {
	data *STATE
}

type DECRYPTION_STATE struct {
	data *STATE
}

func main() {

	type x struct {
		y int
	}

	a := 3

	z := x{a}
	z.y = 7
	fmt.Println(a)

	const key = "YELLOW SUBMARINE"
	keys1 := []byte(key)
	key_bytes := [16]byte{}
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
	block := [16]byte{} //Blocks of ciphertext
	f, err := os.Create("plaintext.txt")
	w := bufio.NewWriter(f)
	for j := 0; j+15 < len(cipher_text_bytes); j += 16 {
		for i := 0; i < 16; i++ {
			block[i] = byte(cipher_text_bytes[i+j])
		}
		state := decrypt(make_decryption_state(cipher_text_bytes[j:j+16]), get_key_schedule(key_bytes))
		plain_text_bytes := unpackState(*state.data)
		fmt.Println(string(plain_text_bytes[:]))
		w.WriteString(string(plain_text_bytes[:]))
	}
	w.Flush()
}

func decrypt(state DECRYPTION_STATE, key_schedule [44]uint32) DECRYPTION_STATE {

	// dec_state := DECRYPTION_STATE{&state}

	//Round 0
	round := 0
	round_key := make_key_matrix(key_schedule[40-round*4 : 44-round*4])
	state.data.add_round_key(round_key)

	for round := 1; round <= 9; round++ {
		state.data.inv_shift_rows()
		state.lookup()
		key_matrix := make_key_matrix(key_schedule[40-4*round : 44-4*round])
		state.data.add_round_key(key_matrix)
		state.data.inv_mix_columns()
	}

	round = 10
	state.data.inv_shift_rows()
	state.lookup()
	round_key = make_key_matrix(key_schedule[40-4*round : 44-4*round])
	state.data.add_round_key(round_key)
	return state
}
func encrypt(plaintext [16]byte, key_schedule [44]uint32) STATE {

	round_key := make_key_matrix(key_schedule[0:4])
	state := initializeState(plaintext[:])
	enc_state := ENCRYPTION_STATE{&state}
	state.add_round_key(round_key)

	for round := 1; round <= 9; round++ {
		// state.substitute_sbox()
		enc_state.lookup()
		state.shift_rows()
		state.MixColumns()
		key_matrix := make_key_matrix(key_schedule[4*round : 4*round+4])
		state.add_round_key(key_matrix)
	}

	round := 10
	enc_state.lookup()
	state.shift_rows()
	round_key = make_key_matrix(key_schedule[4*round : 4*round+4])
	state.add_round_key(round_key)
	return state
}
func show_block(title string, data [16]byte) {
	fmt.Println(title)
	fmt.Println("")
	fmt.Printf("%x", data)
	fmt.Println("")
}

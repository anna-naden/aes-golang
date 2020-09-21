package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
)

func crypto_challenge() {

	//Encryption key
	const key = "YELLOW SUBMARINE"
	key_bytes := []byte(key)
	key_schedule := get_key_schedule(key_bytes)

	//Ciphertext
	cipher_text_bytes := get_cipher_bytes("crypto-challenge.txt")

	//Decrypt, one block at a time
	plain_text := ""
	for i:=0; i<len(cipher_text_bytes); i +=16 {
		block := cipher_text_bytes[i:i+16]
		state := decrypt(make_decryption_state(block),key_schedule)
		plain_text_bytes := unpackState(state)
		plain_text += string(plain_text_bytes[:])
		}
	fmt.Println(plain_text)

	//Save plaintext to file
	f, err := os.Create("plaintext.txt")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	w.WriteString(plain_text)
	w.Flush()

	fmt.Printf("sbox cpu %d microseconds", int(sbox_cpu)/1000)
}

func main() {
	do_stallings := false
	if do_stallings {
		stallings()
	}
	do_challenge := true
	if do_challenge {
		crypto_challenge()
	}
}
func stallings() {
	key, err := hex.DecodeString("0f1571c947d9e8590cb7add6af7f6798")
	if err != nil {
		panic(err)
	}
	plain_text_bytes, err := hex.DecodeString("0123456789abcdeffedcba9876543210")
	key_schedule := get_key_schedule(key)
	cipher_text_state := encrypt(plain_text_bytes, key_schedule)
	cipher_text_state.show_state("after encryption")
	cipher_text_bytes := unpackState(cipher_text_state)
	fmt.Printf("%x", cipher_text_bytes)
}
func decrypt(state STATE, key_schedule [44]uint32) STATE {

	//Round 0
	round := 0
	round_key := make_key_matrix(key_schedule[40-round*4 : 44-round*4])
	state = state.add_round_key(round_key)

	// Rounds 1-9
	for round := 1; round <= 9; round++ {
		key_matrix := make_key_matrix(key_schedule[40-4*round : 44-4*round])
		state = state.inv_shift_rows().inv_lookup().add_round_key(key_matrix).inv_mix_columns()
	}

	// Round 10
	round = 10
	round_key = make_key_matrix(key_schedule[40-4*round : 44-4*round])
	state = state.inv_shift_rows().inv_lookup().add_round_key(round_key)
	return state
}
func encrypt(plaintext []byte, key_schedule [44]uint32) STATE {

	round_key := make_key_matrix(key_schedule[0:4])
	state := initializeState(plaintext)
	state = state.add_round_key(round_key)

	for round := 1; round <= 9; round++ {
		// state.substitute_sbox()
		key_matrix := make_key_matrix(key_schedule[4*round : 4*round+4])
		state = state.lookup().shift_rows().MixColumns().add_round_key(key_matrix)
	}

	round := 10
	state = state.lookup()
	state = state.shift_rows()
	round_key = make_key_matrix(key_schedule[4*round : 4*round+4])
	state = state.add_round_key(round_key)
	return state
}

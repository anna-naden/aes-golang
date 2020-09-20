package main

/*
#include <pthread.h>
#include <time.h>
#include <stdio.h>

static long long getThreadCpuTimeNs() {
    struct timespec t;
    if (clock_gettime(CLOCK_THREAD_CPUTIME_ID, &t)) {
        perror("clock_gettime");
        return 0;
    }
    return t.tv_sec * 1000000000LL + t.tv_nsec;
}
*/
import "C"

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type STATE [4][4]byte

type SBOX_LOOKER_UPPER interface {
	lookup() STATE
}

var g_cache = [256][256]byte{}
var sbox_cpu = int32(0)

func main() {
	for i:=0; i<256; i++ {
		for j:=0; j<256; j++ {
			g_cache[i][j]=GMul(byte(i),byte(j))
		}
	}
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
func crypto_challenge() {
	const key = "YELLOW SUBMARINE"
	keys1 := []byte(key)
	key_bytes := [16]byte{}
	for i := 0; i < 16; i++ {
		key_bytes[i] = keys1[i]
	}
	cipher_text_bytes := get_cipher_bytes("crypto-challenge.txt")
	block := [16]byte{} //Blocks of ciphertext
	f, err := os.Create("plaintext.txt")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	plain_text_bytes := [16]byte{}
	start := time.Now()
	cpu1 := C.getThreadCpuTimeNs()
	for j := 0; j+15 < len(cipher_text_bytes); j += 16 {
		for i := 0; i < 16; i++ {
			block[i] = byte(cipher_text_bytes[i+j])
		}
		state := decrypt(make_decryption_state(cipher_text_bytes[j:j+16]), get_key_schedule(key_bytes[:]))
		plain_text_bytes = unpackState(state)
		fmt.Println(string(plain_text_bytes[:]))
		w.WriteString(string(plain_text_bytes[:]))
	}
	w.Flush()
	cpu2 := C.getThreadCpuTimeNs()
	fmt.Printf("end-to-end cpu %d ns",(cpu2-cpu1))
	fmt.Println("")
	finish :=time.Since(start)
	fmt.Printf("end-to-end wall time %s\n",finish)
	fmt.Printf("sbox cpu %d", sbox_cpu)
	fmt.Println(string(plain_text_bytes[:]))
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
func show_block(title string, data [16]byte) {
	fmt.Println(title)
	fmt.Println("")
	fmt.Printf("%x", data)
	fmt.Println("")
}
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

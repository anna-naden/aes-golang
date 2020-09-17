package main

import (
	// "encoding/hex"
	"fmt"
	"io/ioutil"
	"encoding/base64"
)

func main() {
	// myBytes := [5]uint8 {
	// 	0xff,0x02,0x20,0x00,0xc5,
	// }
	// showBits(myBytes[1])
	// fmt.Printf("%x",sub_bytes(inv_sub_bytes(myBytes[:])))


	// // Test mixcolumns

	// state := [4][4]uint8 {
	// 	{0x87,0xf2,0x4d,0x97},
	// 	{0x6e,0x4c,0x90,0xec},
	// 	{0x46,0xe7,0x4a,0xc3},
	// 	{0xa6,0x8c,0xd8,0x95},
	// }
	// state = [4][4]uint8 {
	// 	{0xab,0x8b,0x89,0x35},
	// 	{0x40,0x7f,0xf1,0x05},
	// 	{0xf0,0xfc,0x18,0x3f},
	// 	{0xc4,0xe4,0x4e,0x2f},
	// }
	// show_state("before mix", state)
	// show_state("after mix", MixColumns(state))
	// show_state("revert", inv_mix_columns(MixColumns(state)))

	// Test inv_shift_rows

	// state := [4][4]uint8 {
	// 	{0xab,0x8b,0x89,0x35},
	// 	{0x40,0x7f,0xf1,0x05},
	// 	{0xf0,0xfc,0x18,0x3f},
	// 	{0xc4,0xe4,0x4e,0x2f},
	// }
	// show_state("before sh/inv", state)
	// show_state("after sh/inv", inv_shift_rows(shift_rows(state)))

	//Stallings case

	// const hex_key = "0f1571c947d9e8590cb7add6af7f6798"
	// bkey, err := hex.DecodeString(hex_key)
	// if err != nil {
	// 	panic(err)
	// }
	// bkey2 := [16]uint8{}
	// for i := 0; i < 16; i++ {
	// 	bkey2[i] = bkey[i]
	// }
	// keys := key_expand(bkey2)
	// hex_plain := "0123456789abcdeffedcba9876543210"
	// plaintext, err := hex.DecodeString(hex_plain)
	// if err != nil {
	// 	panic(err)
	// }
	// p2:=[16]uint8{}
	// for i:=0;i<16;i++ {
	// 	p2[i]= plaintext[i] 
	// }
	const key = "YELLOW SUBMARINE"
	keys1 := []byte(key)
	keys := [16]uint8{}
	for i:=0;i<16;i++ {
		keys[i]=keys1[i]
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
	p2:=[16]uint8{}
	fmt.Println(len(bytes1))
	for j:=0; j+15<2880; j += 16 {
		for i:=0; i<16; i++ {
			p2[i] = uint8(bytes1[i+j])
		}
		state := decrypt(bytesToMatrix(p2), key_expand(keys))
		bytes := matrixToBytes(state)
		// fmt.Println(j)
		fmt.Println(string(bytes[:]))
}
	// state := encrypt(p2, key_expand(keys))
	// show_state(s)
	// x:=uint8(0x57)
	// y:=uint8(0x83)
	// fmt.Printf("% x",key_expand(bkey2))
	// for i:=1;i<=11; i++ {
	// 	fmt.Printf("\ni %d rc %x",i,round_const(i))
	// }
	// fmt.Printf("---%x",subWord(rotWord(0xaf7f6798)))
}

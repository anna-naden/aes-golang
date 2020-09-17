package main

import "encoding/binary"
import "fmt"

func add_round_key(key [4][4]uint8, state [4][4]uint8) [4][4]uint8 {
	out := [4][4]uint8{}
	for row:=0; row<4;row++ {
		for col:=0; col<4; col++ {
			out[row][col] = state[row][col] ^ key[row][col]
		}
	}
	return out
}
func bytesToMatrix(theBytes [16]uint8) [4][4]uint8 {
	out := [4][4]uint8{}
	j:=0
	for row := 0; row<4; row++ {
		for col :=0; col<4; col++ {
			out[col][row] = theBytes[j]
			j++
		}
	}
	return out
}
func bytesToWord(ba [4]uint8) uint32 {
	var value uint32
	value |= uint32(ba[0])<<24
	value |= uint32(ba[1]) << 16
	value |= uint32(ba[2]) << 8
	value |= uint32(ba[3])
	return value
}
func decrypt(ciphertext [4][4]uint8, expanded_key [44]uint32) [4][4]uint8 {

	cipher_matrix := ciphertext

	//Round 0
	round := 0
	key_words := get_key_slice(expanded_key,40-round*4)
	key_matrix := wordsToMatrix(key_words)
	state := add_round_key(key_matrix, cipher_matrix)
	for round:=1; round<=9; round++ {

		state = inv_shift_rows(state)
		// if round == 3 {
		// 	show_state("dec rnd 3 before sub by", state)
		// }
		state = inv_sub_bytes_state(state)
		
		key_words := get_key_slice(expanded_key,40-4*round)
		key_matrix := wordsToMatrix(key_words)
		state = add_round_key(key_matrix, state)

		state = inv_mix_columns(state)
	}

	round = 10
	state = inv_shift_rows(state)
	state = inv_sub_bytes_state(state)

	key_words = get_key_slice(expanded_key,40-round*4)
	key_matrix = wordsToMatrix(key_words)
	state = add_round_key(key_matrix, state)
	// show_state("dec rnd 10 after add rnd key", state)

return state
}
func encrypt(plaintext [16]uint8, expanded_key [44]uint32) [4][4]uint8 {
	key_words := get_key_slice(expanded_key,0)
	key_matrix := wordsToMatrix(key_words)
	plain_matrix := bytesToMatrix(plaintext)
	// show_state("encrypt rnd 0 before add key", plain_matrix)
	state := add_round_key(key_matrix,plain_matrix)
	box := sbox()
	for round:=1; round<=9; round++ {
		for row:=0; row <4; row++ {
			for col:=0; col<4; col++ {
				b := state[row][col]
				j := int(b & 0xf)
				i := int((b & 0xf0) >> 4)
				state[row][col] = box[i][j]
			}
		}
		state = shift_rows(state)
		state = MixColumns(state)
		key_words := get_key_slice(expanded_key,4*round)
		key_matrix := wordsToMatrix(key_words)
		state = add_round_key(key_matrix, state)
	}
	round:=10
	for row:=0; row <4; row++ {
		for col:=0; col<4; col++ {
			b := state[row][col]
			j := int(b & 0xf)
			i := int((b & 0xf0) >> 4)
			state[row][col] = box[i][j]
		}
	}
	state = shift_rows(state)

	key_words = get_key_slice(expanded_key,4*round)
	key_matrix = wordsToMatrix(key_words)
	state = add_round_key(key_matrix, state)
	// show_state("after adding round key", state)
return state
}
func g_aes(word uint32, j int) uint32 {
	word2 := rotWord(word)
	word2 = subWord(word2)
	rc := round_const(j)
	// fmt.Printf("\ng_aes j %d rc %x",j,rc)
	// fmt.Printf("g_aes word2 %x",word2)
	// fmt.Printf("g_aes_xor %x", word2^rc)
	return(word2^rc)
}
func get_key_slice(key_words [44]uint32, i int) [4]uint32 {
	out := [4]uint32{}
	for j:=i; j<i+4;j++ {
		out[j-i] = key_words[j]
	}
	return out
}
func inv_mix_columns(s [4][4]uint8) [4][4]uint8 { 
	// The first index is the row
    ss:=[4][4]uint8{}
    for  c := 0; c < 4; c++ {
		ss[0][c] = 
			GMul(0x0e, s[0][c])^ 
			GMul(0x0b, s[1][c])^
			GMul(0x0d, s[2][c])^
			GMul(0x09, s[3][c])
		ss[1][c] = 
			GMul(0x09, s[0][c])^
			GMul(0x0e, s[1][c])^
			GMul(0x0b, s[2][c])^
			GMul(0x0d, s[3][c])
		ss[2][c] =
			GMul(0x0d, s[0][c])^
			GMul(0x09, s[1][c])^
			GMul(0x0e, s[2][c])^
			GMul(0x0b, s[3][c])
		ss[3][c] =
			GMul(0x0b, s[0][c])^
			GMul(0x0d, s[1][c])^
			GMul(0x09, s[2][c])^
			GMul(0x0e, s[3][c])
    }
    return ss
}
func inv_sub_bytes(bytes []uint8) []uint8 {
	retval := []uint8{}
	box := inv_sbox()
	for _, b := range bytes {
		col := b & 0xf
		row := (b & 0xf0) >> 4
		lookup :=    box[row][col]
		retval = append(retval, lookup)
	}
	return retval

}
func inv_sub_bytes_state(state [4][4]uint8) [4][4]uint8 {
	retval := [4][4]uint8{}
	box := inv_sbox()
	for row :=0; row<4; row++ {
		for col :=0; col<4; col++ {
			b := state[row][col]
			i := (b & 0xf0) >> 4
			j := b & 0xf
			lookup :=box[i][j]
			retval[row][col] = lookup
		}
	}
	return retval
}

func inv_shift_rows(input [4][4]uint8) [4][4]uint8 {
	output := [4][4]uint8 {}
	output[0]=input[0]
	
	in_row1 := input[1]
	out_row1 := [4]uint8{}
	for i :=0; i<4; i++ {
		j:=i-1
		if j<0 {
			j +=4
		}
		out_row1[i] = in_row1[j]
	}
	output[1] = out_row1

	in_row2 := input[2]
	out_row2 := [4]uint8{}
	for i:=0; i<4; i++ {
		j:=i-2
		if j<0 {
			j +=4
		}
		out_row2[i]=in_row2[j]
	}
	output[2]=out_row2
	
	in_row3 := input[3]
	out_row3 := [4]uint8{}
	for i:=0; i<4; i++ {
		j:=i-3
		if j<0 {
			j +=4
		}
		out_row3[i] = in_row3[j]
	}
	output[3]=out_row3

	return output
}
func inv_sbox() [16][16]uint8 {
	box := [256]uint8{
		0x52, 0x09, 0x6a, 0xd5, 0x30, 0x36, 0xa5, 0x38, 0xbf, 0x40, 0xa3, 0x9e, 0x81, 0xf3, 0xd7, 0xfb,
		0x7c, 0xe3, 0x39, 0x82, 0x9b, 0x2f, 0xff, 0x87, 0x34, 0x8e, 0x43, 0x44, 0xc4, 0xde, 0xe9, 0xcb,
		0x54, 0x7b, 0x94, 0x32, 0xa6, 0xc2, 0x23, 0x3d, 0xee, 0x4c, 0x95, 0x0b, 0x42, 0xfa, 0xc3, 0x4e,
		0x08, 0x2e, 0xa1, 0x66, 0x28, 0xd9, 0x24, 0xb2, 0x76, 0x5b, 0xa2, 0x49, 0x6d, 0x8b, 0xd1, 0x25,
		0x72, 0xf8, 0xf6, 0x64, 0x86, 0x68, 0x98, 0x16, 0xd4, 0xa4, 0x5c, 0xcc, 0x5d, 0x65, 0xb6, 0x92,
		0x6c, 0x70, 0x48, 0x50, 0xfd, 0xed, 0xb9, 0xda, 0x5e, 0x15, 0x46, 0x57, 0xa7, 0x8d, 0x9d, 0x84,
		0x90, 0xd8, 0xab, 0x00, 0x8c, 0xbc, 0xd3, 0x0a, 0xf7, 0xe4, 0x58, 0x05, 0xb8, 0xb3, 0x45, 0x06,
		0xd0, 0x2c, 0x1e, 0x8f, 0xca, 0x3f, 0x0f, 0x02, 0xc1, 0xaf, 0xbd, 0x03, 0x01, 0x13, 0x8a, 0x6b,
		0x3a, 0x91, 0x11, 0x41, 0x4f, 0x67, 0xdc, 0xea, 0x97, 0xf2, 0xcf, 0xce, 0xf0, 0xb4, 0xe6, 0x73,
		0x96, 0xac, 0x74, 0x22, 0xe7, 0xad, 0x35, 0x85, 0xe2, 0xf9, 0x37, 0xe8, 0x1c, 0x75, 0xdf, 0x6e,
		0x47, 0xf1, 0x1a, 0x71, 0x1d, 0x29, 0xc5, 0x89, 0x6f, 0xb7, 0x62, 0x0e, 0xaa, 0x18, 0xbe, 0x1b,
		0xfc, 0x56, 0x3e, 0x4b, 0xc6, 0xd2, 0x79, 0x20, 0x9a, 0xdb, 0xc0, 0xfe, 0x78, 0xcd, 0x5a, 0xf4,
		0x1f, 0xdd, 0xa8, 0x33, 0x88, 0x07, 0xc7, 0x31, 0xb1, 0x12, 0x10, 0x59, 0x27, 0x80, 0xec, 0x5f,
		0x60, 0x51, 0x7f, 0xa9, 0x19, 0xb5, 0x4a, 0x0d, 0x2d, 0xe5, 0x7a, 0x9f, 0x93, 0xc9, 0x9c, 0xef,
		0xa0, 0xe0, 0x3b, 0x4d, 0xae, 0x2a, 0xf5, 0xb0, 0xc8, 0xeb, 0xbb, 0x3c, 0x83, 0x53, 0x99, 0x61,
		0x17, 0x2b, 0x04, 0x7e, 0xba, 0x77, 0xd6, 0x26, 0xe1, 0x69, 0x14, 0x63, 0x55, 0x21, 0x0c, 0x7d}
	box2 := [16][16]uint8{}
	for col := 0; col < 16; col++ {
		for row := 0; row < 16; row++ {
			box2[row][col] = box[row*16+col]
		}
	}
	return box2
}

func key_expand(key [16]uint8 ) [44]uint32 {
	words := [44]uint32{}
	for i:=0; i<4;i++ {
		words[i] = binary.BigEndian.Uint32(key[4*i:4*i+4])
		// fmt.Printf("% x", words[i])
	}
	for i:=1;i<11;i++ {
		g := g_aes(words[4*i-1],i)
		words[4*i] = words[4*i-4]^g
		words[4*i+1] = words[4*i-3]^words[4*i]
		words[4*i+2] = words[4*i-2]^words[4*i+1]
		words[4*i+3] = words[4*i-1]^words[4*i+2]
	}
	// fmt.Printf("% x",words)
	retval := [44]uint32{}
	for i:= 0;i<44;i++ {
		retval[i]=words[i]
	}
	return retval
}

func GMul(a uint8, b uint8) uint8 {
	 // Galois Field (256) Multiplication of two Bytes
	p:=uint8(0)
    for counter := 0; counter < 8; counter++ {
        if ((b & 1) != 0) {
            p ^= a
        }

        hi_bit_set := (a & 0x80) != 0
        a <<= 1
        if (hi_bit_set) {
            a ^= 0x1B; // x^8 + x^4 + x^3 + x + 1 */
        }
        b >>= 1;
    }

    return p;
}
func matrixToBytes(matrix [4][4]uint8) [16]uint8 {
	out := [16]uint8{}
	j:=0
	for row:=0;row<4;row++ {
		for col:=0; col<4; col++ {
			out[j] = matrix[col][row]
			j++
		}
	}
	return out
}
func MixColumns(s [4][4]uint8) [4][4]uint8 { 
	// 's' is the main State matrix, 'ss' is a temp matrix of the same dimensions as 's'.
	// The first index is the row
    ss:=[4][4]uint8{}
    for  c := 0; c < 4; c++ {
        ss[0][c] = (GMul(0x02, s[0] [c]) ^ GMul(0x03, s[1] [c]) ^ s[2] [c] ^ s[3][c])
        ss[1][c] = (s[0][c] ^ GMul(0x02, s[1][c]) ^ GMul(0x03, s[2][c]) ^ s[3][c]);
        ss[2][c] = (s[0][c] ^ s[1][c] ^ GMul(0x02, s[2][c]) ^ GMul(0x03, s[3] [c]));
        ss[3][c] = (GMul(0x03, s[0][c]) ^ s[1][c] ^ s[2][c] ^ GMul(0x02, s[3] [c]));
    }

    return ss
}
func round_const(j int) uint32 {
	if j<9 {
		return 1<<(23+j)
	} else if j==9 {
		return 0x1b000000 
	} else {
		return 0x36000000
	}
}
func rotWord(word uint32) uint32 {
	ba := wordToBytes(word)
	b0:=ba[0]
	var ba2 [4]uint8
	ba2 = [4]uint8{}
	for i:=0;i<3;i++ {
		ba2[i]=ba[i+1]
	}
	ba2[3]=b0
	return(bytesToWord(ba2))
}
func sbox() [16][16]uint8 {
	box := [256]uint8{
		0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76,
		0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0,
		0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15,
		0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75,
		0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84,
		0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf,
		0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8,
		0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2,
		0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73,
		0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb,
		0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79,
		0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08,
		0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a,
		0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e,
		0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf,
		0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16}
	box2 := [16][16]uint8{}
	for col := 0; col < 16; col++ {
		for row := 0; row < 16; row++ {
			box2[row][col] = box[row*16+col]
		}
	}
	return box2
}
func shift_rows(input [4][4]uint8) [4][4]uint8 {
	output := [4][4]uint8 {}
	output[0]=input[0]
	
	in_row1 := input[1]
	out_row1 := [4]uint8{}
	out_row1[0] = in_row1[1]
	out_row1[1] = in_row1[2]
	out_row1[2] = in_row1[3]
	out_row1[3] = in_row1[0]
	output[1] = out_row1

	in_row2 := input[2]
	out_row2 := [4]uint8{}
	out_row2[0] = in_row2[2]
	out_row2[1] = in_row2[3]
	out_row2[2] = in_row2[0]
	out_row2[3] = in_row2[1]
	output[2] = out_row2

	in_row3 := input[3]
	out_row3 := [4]uint8{}
	out_row3[0]=in_row3[3]
	out_row3[1]=in_row3[0]
	out_row3[2]=in_row3[1]
	out_row3[3]=in_row3[2]
	output[3]=out_row3

	return output
}
func show_state(title string,s [4][4]uint8) {
	fmt.Println(title)
	for row:=0;row<4;row++ {
		for col:=0;col<4;col++ {
			fmt.Printf("%x ",s[row][col])
		}
		fmt.Println("")
	}
	fmt.Println("")
}	
func sub_bytes(bytes []uint8) []uint8 {
	retval := []uint8{}
	for _, b := range bytes {
		col := b & 0xf
		row := (b & 0xf0) >> 4
		box := sbox()
		retval = append(retval, box[row][col])
	}
	return retval
}

func sub_bytes16(bytes [16]uint8) [16]uint8 {
	retval := [16]uint8{}
	box := sbox()
	for i:=0; i<16; i++ {
		b := bytes[i]
		col := b & 0xf
		row := (b & 0xf0) >> 4
		retval[i] = box[row][col]
	}
	return retval
}

func subWord(word uint32) uint32 {
	myBytes := wordToBytes(word)
	ba := sub_bytes(myBytes[:])
	ba2 := [4]uint8{ba[0],ba[1],ba[2],ba[3]}
	return bytesToWord(ba2)
	
}
func wordsToMatrix(words [4]uint32) [4][4]uint8 {
	out:= [4][4]uint8{}
	for row:=0;row<4;row++ {
		for col:=0; col<4; col++ {
			word := words[row]
			ishift := 24-col*8
			theByte := uint8((word&(0xff<<ishift))>>ishift)
			out[col][row] = theByte
		}
	}
	return out
}
func wordToBytes(word uint32) [4]uint8 {
	myBytes := [4]uint8 {
		uint8((word&0xff000000)>>24),
		uint8((word&0xff0000)>>16), 
		uint8((word&0xff00)>>8), 
		uint8(word&0xff)}
		return myBytes
}
//UNUSED ATTEMPT TO COMPUTE MULTIPLICATIVE INVERSE
func mult_inv(inByte uint8) uint8 {
	inbits := [8]uint8{}
	for i := 0; i < 8; i++ {
		bi := (inByte & (1 << i)) >> i
		inbits[i]=bi
	}
	const d=0x05
	outByte:=uint8(0)
	showBits(inByte)
	for i:=0;i<8;i++ {
		outbit := inbits[(i+2)%8]^inbits[(i+5)%8]^inbits[(i+7)%8]^d
		if outbit == 1 {
			outByte += 1<<i
		}
	}
	showBits(outByte)
	return outByte
}



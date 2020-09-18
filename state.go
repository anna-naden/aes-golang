package main

import "fmt"

func (state *STATE) add_round_key(key [4][4]uint8) {
	out := STATE{}
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			out[row][col] = state[row][col] ^ key[row][col]
		}
	}
	*state = out
}

func (s *STATE) inv_mix_columns() {
	// The first index is the row
	ss := STATE{}
	for c := 0; c < 4; c++ {
		ss[0][c] =
			GMul(0x0e, s[0][c]) ^
				GMul(0x0b, s[1][c]) ^
				GMul(0x0d, s[2][c]) ^
				GMul(0x09, s[3][c])
		ss[1][c] =
			GMul(0x09, s[0][c]) ^
				GMul(0x0e, s[1][c]) ^
				GMul(0x0b, s[2][c]) ^
				GMul(0x0d, s[3][c])
		ss[2][c] =
			GMul(0x0d, s[0][c]) ^
				GMul(0x09, s[1][c]) ^
				GMul(0x0e, s[2][c]) ^
				GMul(0x0b, s[3][c])
		ss[3][c] =
			GMul(0x0b, s[0][c]) ^
				GMul(0x0d, s[1][c]) ^
				GMul(0x09, s[2][c]) ^
				GMul(0x0e, s[3][c])
	}
	*s = ss
	return
}

func (input *STATE) inv_shift_rows() {
	output := STATE{}
	output[0] = input[0]

	in_row1 := input[1]
	out_row1 := [4]uint8{}
	for i := 0; i < 4; i++ {
		j := i - 1
		if j < 0 {
			j += 4
		}
		out_row1[i] = in_row1[j]
	}
	output[1] = out_row1

	in_row2 := input[2]
	out_row2 := [4]uint8{}
	for i := 0; i < 4; i++ {
		j := i - 2
		if j < 0 {
			j += 4
		}
		out_row2[i] = in_row2[j]
	}
	output[2] = out_row2

	in_row3 := input[3]
	out_row3 := [4]uint8{}
	for i := 0; i < 4; i++ {
		j := i - 3
		if j < 0 {
			j += 4
		}
		out_row3[i] = in_row3[j]
	}
	output[3] = out_row3

	*input = output
}

func (s *STATE) MixColumns() {
	ss := [4][4]uint8{}
	for c := 0; c < 4; c++ {
		ss[0][c] = (GMul(0x02, s[0][c]) ^ GMul(0x03, s[1][c]) ^ s[2][c] ^ s[3][c])
		ss[1][c] = (s[0][c] ^ GMul(0x02, s[1][c]) ^ GMul(0x03, s[2][c]) ^ s[3][c])
		ss[2][c] = (s[0][c] ^ s[1][c] ^ GMul(0x02, s[2][c]) ^ GMul(0x03, s[3][c]))
		ss[3][c] = (GMul(0x03, s[0][c]) ^ s[1][c] ^ s[2][c] ^ GMul(0x02, s[3][c]))
	}

	*s = ss
}

func (input *STATE) shift_rows() {
	output := STATE{}
	output[0] = input[0]

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
	out_row3[0] = in_row3[3]
	out_row3[1] = in_row3[0]
	out_row3[2] = in_row3[1]
	out_row3[3] = in_row3[2]
	output[3] = out_row3

	*input = output
}
func (s STATE) show_state(title string) {
	fmt.Println(title)
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			fmt.Printf("%x ", s[row][col])
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
		box := get_sbox()
		retval = append(retval, box[row][col])
	}
	return retval
}

func (dec_state *DECRYPTION_STATE) lookup() {
	box := get_inv_sbox()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			b := &(dec_state.state[row][col])
			i := int((*b & 0xf0) >> 4)
			j := int(*b & 0xf)
			*b = box[i][j]
		}
	}
}
func (enc_state *ENCRYPTION_STATE) lookup() {
	box := get_sbox()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			b := &(enc_state.state[row][col])
			j := int(*b & 0xf)
			i := int((*b & 0xf0) >> 4)
			*b = box[i][j]
		}
	}

}


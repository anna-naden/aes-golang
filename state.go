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
import "fmt"

type STATE [4][4]byte

var sbox_cpu = int32(0)

var cpu_used = uint32(0)
var g_cache_fetched = false
var g_cache = [256][256]byte{}
func (s STATE) inv_mix_columns() STATE{
	if !g_cache_fetched {
		for i:=0; i<256; i++ {
			for j:=0; j<256; j++ {
				g_cache[i][j]=galois_multiply(byte(i),byte(j))
			}
		}
		g_cache_fetched = true
	}
	f:= func(s STATE)STATE{
		var mixer = [4][4]byte {
			{0x0e, 0x0b, 0x0d, 0x09},
			{0x09, 0x0e, 0x0b, 0x0d},
			{0x0d, 0x09, 0x0e, 0x0b},
			{0x0b, 0x0d, 0x09, 0x0e},
		}
		cpu_start := C.getThreadCpuTimeNs()
		// The first index is the row
		ss := STATE{}
		for row := 0; row<4; row ++ {
			for col :=0; col<4; col++ {
				matrix_element := byte(0)
				for i:=0; i<4; i++ {
					x:=mixer[row][i]
					y := s[i][col]
					matrix_element ^= g_cache[x][y]
				}
				ss[row][col] = matrix_element
			}
		}
		cpu_used += uint32(C.getThreadCpuTimeNs() - cpu_start)
		return ss
	}
	return f(s)
}

func (state STATE) add_round_key(key [4][4]byte) STATE {
	out := STATE{}
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			out[row][col] = state[row][col] ^ key[row][col]
		}
	}
	return out
}

func (input STATE) inv_shift_rows() STATE {
	output := STATE{}
	output[0] = input[0]

	in_row1 := input[1]
	out_row1 := [4]byte{}
	for i := 0; i < 4; i++ {
		j := i - 1
		if j < 0 {
			j += 4
		}
		out_row1[i] = in_row1[j]
	}
	output[1] = out_row1

	in_row2 := input[2]
	out_row2 := [4]byte{}
	for i := 0; i < 4; i++ {
		j := i - 2
		if j < 0 {
			j += 4
		}
		out_row2[i] = in_row2[j]
	}
	output[2] = out_row2

	in_row3 := input[3]
	out_row3 := [4]byte{}
	for i := 0; i < 4; i++ {
		j := i - 3
		if j < 0 {
			j += 4
		}
		out_row3[i] = in_row3[j]
	}
	output[3] = out_row3

	return output
}

func (s STATE) MixColumns() STATE {
	ss := [4][4]byte{}
	for c := 0; c < 4; c++ {
		ss[0][c] = (galois_multiply(0x02, s[0][c]) ^ galois_multiply(0x03, s[1][c]) ^ s[2][c] ^ s[3][c])
		ss[1][c] = (s[0][c] ^ galois_multiply(0x02, s[1][c]) ^ galois_multiply(0x03, s[2][c]) ^ s[3][c])
		ss[2][c] = (s[0][c] ^ s[1][c] ^ galois_multiply(0x02, s[2][c]) ^ galois_multiply(0x03, s[3][c]))
		ss[3][c] = (galois_multiply(0x03, s[0][c]) ^ s[1][c] ^ s[2][c] ^ galois_multiply(0x02, s[3][c]))
	}

	return ss
}

func (input STATE) shift_rows() STATE {
	output := STATE{}
	output[0] = input[0]

	in_row1 := input[1]
	out_row1 := [4]byte{}
	out_row1[0] = in_row1[1]
	out_row1[1] = in_row1[2]
	out_row1[2] = in_row1[3]
	out_row1[3] = in_row1[0]
	output[1] = out_row1

	in_row2 := input[2]
	out_row2 := [4]byte{}
	out_row2[0] = in_row2[2]
	out_row2[1] = in_row2[3]
	out_row2[2] = in_row2[0]
	out_row2[3] = in_row2[1]
	output[2] = out_row2

	in_row3 := input[3]
	out_row3 := [4]byte{}
	out_row3[0] = in_row3[3]
	out_row3[1] = in_row3[0]
	out_row3[2] = in_row3[1]
	out_row3[3] = in_row3[2]
	output[3] = out_row3

	return output
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

var s_box_fetched = false
var s_box = [16][16]byte{}

func sub_bytes(bytes []byte) []byte {
	f:= func(bytes []byte) []byte {
		if !s_box_fetched {
			s_box = get_sbox()
			s_box_fetched=true
		}
		retval := []byte{}
		for _, b := range bytes {
			col := b & 0xf
			row := (b & 0xf0) >> 4
			retval = append(retval, s_box[row][col])
		}
		return retval
	}
	return f(bytes)
}

var inv_sbox_fetched = false
var inv_sbox = [16][16]byte {}

func (dec_state STATE) inv_lookup() STATE {
	cpu1 := C.getThreadCpuTimeNs()
	f:= func(state STATE) STATE {
		if !inv_sbox_fetched {
			inv_sbox = get_inv_sbox()
			inv_sbox_fetched = true
		}
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				b := &(dec_state[row][col])
				i := int((*b & 0xf0) >> 4)
				j := int(*b & 0xf)
				*b = inv_sbox[i][j]
			}
		}
		cpu2 := C.getThreadCpuTimeNs()
		sbox_cpu += int32(cpu2-cpu1)
		return dec_state
	}
	return f(dec_state)
}
func (enc_state *STATE) lookup() STATE {
	box := get_sbox()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			b := &(enc_state[row][col])
			j := int(*b & 0xf)
			i := int((*b & 0xf0) >> 4)
			*b = box[i][j]
		}
	}
	return *enc_state

}

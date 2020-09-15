package main
import "fmt"

func main() {
	var b byte = byte(255)
	var state [16] byte
	state[0]=b
	state2 := [][]byte {
		[]byte{255,255,0,0},
		[]byte{0,0,255,255},
		[]byte{0,0,255,0},
		[]byte{0,0,0,255},
	}
	fmt.Println(state[0:4], state2)

}
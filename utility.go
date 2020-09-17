package main

import "fmt"

func showBits(inByte uint8) {
	fmt.Println("---------------------------")
	for i := 0; i<8; i++ {
		fmt.Println(inByte & (1<<i))
	}
}
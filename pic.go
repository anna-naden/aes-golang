package main
import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	retval := make([][]uint8, dy, dy)
	for row := 0; row < dy; row++ {
		retval[row] = make([]uint8, dx, dx)
		for column:=0;column<dx;column++ {
			retval[row][column]=uint8(row^column)
		}
	}
	return retval
}

func main() {
	pic.Show(Pic)
}

package main
import "fmt"
import "strings"

func main() {
	char_freq1 := char_freqs()
	fmt.Println(char_freq1)
	s := "This is a test"
	fmt.Println(WordCount(s))
}

func char_freqs() map[string] float32 {
	var freq_map = map[string] float32{
        "a": 0.06990291262135923,
        "b": 0.009708737864077669,
        "c": 0.02330097087378641,
        "d": 0.0458252427184466,
        "e": 0.09825242718446602,
        "f": 0.0170873786407767,
        "g": 0.014757281553398059,
        "h": 0.05242718446601942,
        "i": 0.051262135922330095,
        "j": 0.0003883495145631068,
        "k": 0.004271844660194175,
        "l": 0.04932038834951456,
        "m": 0.012815533980582524,
        "n": 0.060582524271844664,
        "o": 0.06485436893203883,
        "p": 0.016699029126213592,
        "q": 0.0007766990291262136,
        "r": 0.047766990291262135,
        "s": 0.04233009708737864,
        "t": 0.06757281553398058,
        "u": 0.018640776699029128,
        "v": 0.012427184466019418,
        "w": 0.02446601941747573,
        "x": 0.0011650485436893205,
        "y": 0.019417475728155338,
        "z": 0.01,
	}
	return freq_map
}
func WordCount(s string) map[string]int {
	retval := make(map[string]int)
	fields := strings.Split(s,"")
	for _, field := range fields {
		_, ok := retval[field]
		if !ok {
			retval[field] = 1
		} else {
			retval[field]++
		}
	}
	return retval
}


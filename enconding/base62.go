package enconding

import (
	"bytes"
	"math"
)

// All characters
const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Convert number to base62
func Encode(n int) string {
	if n == 0 {
		return string(alphabet[0])
	}

	chars := make([]byte, 0)

	length := len(alphabet)

	for n > 0 {
		r := n / length
		remainder := n % length
		chars = append(chars, alphabet[remainder])
		n = r
	}

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}

// convert base62 token to int
func Decode(t string) int {
	n := 0
	idx := 0.0
	chars := []byte(alphabet)

	charsLength := float64(len(chars))
	tokenLength := float64(len(t))

	for _, c := range []byte(t) {
		power := tokenLength - (idx + 1)
		ind := bytes.IndexByte(chars, c)
		n += ind * int(math.Pow(charsLength, power))
		idx++
	}

	return n
}

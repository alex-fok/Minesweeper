package utils

import (
	"math/rand"
	"strings"
)

const idLetters string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const bits = 6               // Number of bits representing an index; Log2(len(idLetters)) + 1
const mask = 1<<bits - 1     // Masking for 6 bit value
const count = int(63 / bits) // Number of indices per 63 bit

// Generate 16 letter ID
func GenerateId() string {
	sb := strings.Builder{}
	sb.Grow(16)
	for idIdx, value, remain := 0, rand.Int63(), count; idIdx < 16; {
		if remain == 0 {
			value, remain = rand.Int63(), count
		}
		if letterIdx := int(value & mask); letterIdx < len(idLetters) {
			sb.WriteByte(idLetters[letterIdx])
			idIdx++
		}
		value >>= bits
		remain--
	}
	return sb.String()
}

package utils

import (
	"math/rand"
	"strings"
)

// Refer to: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const idLetters string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const bits = 6               // Number of bits representing an index; Log2(len(idLetters)) + 1
const mask = 1<<bits - 1     // Masking for 6 bit value
const count = int(63 / bits) // Number of indices per 63 bit

// Generate 16 letter ID
func CreateClientId() string {
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

func CreateRoomId() uint {
	return uint(rand.Uint32() % 10000)
}

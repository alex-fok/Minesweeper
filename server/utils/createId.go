package utils

import "strings"

// Refer to: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const idLetters string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const bits = 6               // Number of bits representing an index; Log2(len(idLetters)) + 1
const mask = 1<<bits - 1     // Masking for 6 bit value
const count = int(63 / bits) // Number of indices per 63 bit

func createId(length int) string {
	sb := strings.Builder{}
	sb.Grow(length)
	random := getRand()
	for idIdx, value, remain := 0, random.Int63(), count; idIdx < length; {
		if remain == 0 {
			value, remain = random.Int63(), count
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

// Generate 16 letter ID
func CreateClientId() string {
	return createId(16)
}

// Generate 32 letter ID
func CreateInvitationId() string {
	return createId(32)
}

// Number Only
func CreateRoomId() uint {
	random := getRand()
	return uint(random.Uint32() % 10000)
}

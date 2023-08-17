package utils

import (
	"math/rand"
	"time"
)

func getRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GetRandArray(size int) *[]int {
	array := make([]int, size)
	random := getRand()
	for i := range array {
		array[i] = i
	}
	for i := size; i > 0; i-- {
		target := random.Intn(int(i))
		array[i-1], array[target] = array[target], array[i-1]
	}
	return &array
}

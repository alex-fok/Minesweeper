package boardhelper

import (
	// "fmt"
	"math/rand"
)

const (
	BLANK int = 0
	BOMB  int = 1
	NUM   int = 2
)

type Block struct {
	bType int
	val   int
}

func getBombLoc(num int, size int) [][]int {
	max := size * size
	randArr := make([]int, max)

	// Assign init values
	for i := range randArr {
		randArr[i] = i
	}
	// Fisher-Yates Shuffle
	for i := max; i > 0; i-- {
		target := rand.Intn(int(i))
		randArr[i-1], randArr[target] = randArr[target], randArr[i-1]
	}
	// Trasform to x, y location slice
	loc := make([][]int, num)
	for i, v := range randArr[:num] {
		loc[i] = append(loc[i], v/size, v%size)
	}
	return loc
}

func getValidNgh(loc []int, size int) [][]int {
	isValid := func(x int, y int) bool {
		return x >= 0 && x < size && y >= 0 && y < size
	}
	slice := make([][]int, 0, 8)
	row, col := loc[0], loc[1]
	for i := row - 1; i < row+2; i++ {
		for j := col - 1; j < col+2; j++ {
			if i == row && j == col {
				continue
			}
			if isValid(i, j) {
				slice = append(slice, []int{i, j})
			}
		}
	}
	return slice
}

func GetBoard(num int, size int) [][]Block {
	board := make([][]Block, size)

	// Init board
	for i := range board {
		board[i] = make([]Block, size)
		for j := range board[i] {
			board[i][j] = Block{bType: BLANK, val: 0}
		}
	}
	// Place bomb blocks
	bombLoc := getBombLoc(num, size)
	for _, v := range bombLoc {
		row, col := v[0], v[1]
		board[row][col].bType = BOMB
		board[row][col].val = 0
	}
	// Place number blocks
	for _, v := range bombLoc {
		placeables := getValidNgh(v, size)
		for _, v := range placeables {
			row, col := v[0], v[1]
			switch board[row][col].bType {
			case NUM:
				board[row][col].val += 1
			case BLANK:
				board[row][col].bType = NUM
				board[row][col].val = 1
			default: // BOMB
				continue
			}
		}
	}
	// Print board
	// for _, v := range board {
	// 	fmt.Printf("%v\n", v)
	// }
	return board
}

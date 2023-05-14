package boardhelper

import (
	"math/rand"
)

const (
	BLANK int = 0
	BOMB  int = 1
	NUM   int = 2
)

type Vertex struct {
	X, Y int
}

type Block struct {
	Type int `json:"bType"`
	Val  int `json:"value"`
}

type BlockInfo struct {
	X int `json:"x"`
	Y int `json:"y"`
	Block
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

func getNeighbors(loc []int, size int) [][]int {
	isValid := func(x int, y int) bool {
		return x >= 0 && x < size && y >= 0 && y < size
	}
	neighbors := make([][]int, 0, 8)
	row, col := loc[0], loc[1]
	for i := row - 1; i < row+2; i++ {
		for j := col - 1; j < col+2; j++ {
			if i == row && j == col {
				continue
			}
			if isValid(i, j) {
				neighbors = append(neighbors, []int{i, j})
			}
		}
	}
	return neighbors
}

func GetRevealables(v *Vertex, b [][]Block) []BlockInfo {
	source := BlockInfo{
		X:     v.X,
		Y:     v.Y,
		Block: b[v.Y][v.X],
	}
	if source.Block.Type != BLANK {
		return []BlockInfo{source}
	}
	// Init isVisited
	size := len(b)
	isVisited := make([][]bool, size)
	visitedRow := make([]bool, size*size)
	for i := range isVisited {
		isVisited[i], visitedRow = visitedRow[:size], visitedRow[size:]
	}
	isVisited[v.Y][v.X] = true

	// Get all revealables
	var curr BlockInfo
	revealables := []BlockInfo{}
	blocks := []BlockInfo{source}
	for {
		if len(blocks) == 0 {
			return revealables
		}
		curr, blocks = blocks[0], blocks[1:]
		revealables = append(revealables, curr)

		// If the block is NUM, continue
		if b[curr.Y][curr.X].Type == NUM {
			continue
		}
		// Else if block is BLANK, find neighbors
		neighbors := getNeighbors([]int{curr.Y, curr.X}, size)
		for _, block := range neighbors {
			y, x := block[0], block[1]
			// Add to queue if neighbor block is not visited && Not BOMB
			if isVisited[y][x] || b[y][x].Type == BOMB {
				continue
			}
			isVisited[y][x] = true
			blocks = append(blocks, BlockInfo{
				X:     x,
				Y:     y,
				Block: b[y][x],
			})
		}
	}
}

func GetBoard(num int, size int) [][]Block {
	board := make([][]Block, size)

	// Init board
	for i := range board {
		board[i] = make([]Block, size)
		for j := range board[i] {
			board[i][j] = Block{Type: BLANK, Val: 0}
		}
	}
	// Place bomb blocks
	bombLoc := getBombLoc(num, size)
	for _, v := range bombLoc {
		row, col := v[0], v[1]
		board[row][col].Type = BOMB
		board[row][col].Val = 0
	}
	// Place number blocks
	for _, v := range bombLoc {
		placeables := getNeighbors(v, size)
		for _, v := range placeables {
			row, col := v[0], v[1]
			switch board[row][col].Type {
			case NUM:
				board[row][col].Val += 1
			case BLANK:
				board[row][col].Type = NUM
				board[row][col].Val = 1
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

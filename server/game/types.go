package game

import "minesweeper/types"

type ClientId = types.ClientId

type Block struct {
	Type      int `json:"bType"`
	Val       int `json:"value"`
	Visited   bool
	VisitedBy ClientId `json:"visitedBy"`
}

type BlockInfo struct {
	X int `json:"x"`
	Y int `json:"y"`
	*Block
}

type Vertex struct {
	X, Y int
}

type Counter struct {
	Score     map[ClientId]uint `json:"score"`
	BombsLeft uint              `json:"bombsLeft"`
}

type Turn struct {
	Count uint
	Curr  ClientId
	Next  ClientId
}

type Game struct {
	Size     uint
	Bomb     uint
	Counter  Counter
	GameMode string
	Turn     Turn
	Board    [][]*Block
	Winner   ClientId
}

package game

import (
	"math/rand"
)

const DEFAULT_SIZE = 26
const DEFAULT_BOMB_COUNT = 100

func newGame() *Game {
	return &Game{
		Counter: Counter{},
		Turn: Turn{
			Curr: "",
			Next: "",
		},
		Board:  GetBoard(DEFAULT_SIZE, DEFAULT_BOMB_COUNT),
		Winner: "",
	}
}

func (g *Game) initCounter() {
	g.Counter.BombsLeft = DEFAULT_BOMB_COUNT
	g.Counter.Score = make(map[ClientId]uint)
	g.Counter.Score[g.Turn.Curr] = 0
	g.Counter.Score[g.Turn.Next] = 0
}

func (g *Game) getCounter() Counter {
	return g.Counter
}

func (g *Game) getTurn() Turn {
	return g.Turn
}

func (g *Game) getBoard() [][]*Block {
	return g.Board
}

func (g *Game) getWinner() ClientId {
	return g.Winner
}

func (g *Game) reveal(v *Vertex) []*BlockInfo {
	vertices := GetRevealableVertices(v, g.Board)
	curr := g.getTurn().Curr
	revealables := make([]*BlockInfo, len(vertices))

	for i, v := range vertices {
		g.Board[v.Y][v.X].Visited = true
		g.Board[v.Y][v.X].VisitedBy = curr
		revealables[i] = &BlockInfo{
			X:     v.X,
			Y:     v.Y,
			Block: g.Board[v.Y][v.X],
		}
	}
	return revealables
}

func (g *Game) getVisibleBlocks() []BlockInfo {
	s := []BlockInfo{}
	for i := range g.Board {
		for j := range g.Board {
			if g.Board[i][j].Visited {
				s = append(s, BlockInfo{
					X:     j,
					Y:     i,
					Block: g.Board[i][j],
				})
			}
		}
	}
	return s
}

func (g *Game) assignTurn(cId ClientId) (ClientId, ClientId) {
	isEmpty := g.Turn.Curr == "" && g.Turn.Next == ""

	if isEmpty {
		if rand.Intn(2) == 0 {
			g.Turn.Curr = cId
		} else {
			g.Turn.Next = cId
		}
	} else if g.Turn.Curr == "" {
		g.Turn.Curr = cId
	} else if g.Turn.Next == "" { // Not 'else'. Could be more than 3 clients in a room
		g.Turn.Next = cId
	}
	return g.Turn.Curr, g.Turn.Next
}

func (g *Game) unassignTurn(cId ClientId) (ClientId, ClientId) {
	if g.Turn.Curr == cId {
		g.Turn.Curr = ""
	} else if g.Turn.Next == cId {
		g.Turn.Next = ""
	}
	return g.Turn.Curr, g.Turn.Next
}

func (g *Game) advanceTurn() Turn {
	g.Turn.Count++
	g.Turn.Curr, g.Turn.Next = g.Turn.Next, g.Turn.Curr
	return g.Turn
}

func (g *Game) scoreCurrPlayer() (Counter, bool) {
	g.Counter.BombsLeft--
	g.Counter.Score[g.Turn.Curr]++

	isWon := g.Counter.Score[g.Turn.Curr] > DEFAULT_BOMB_COUNT/2
	if isWon {
		g.Winner = g.Turn.Curr
	}
	return g.Counter, isWon
}

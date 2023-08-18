package game

import (
	"log"
	"minesweeper/utils"
)

// const DEFAULT_SIZE = 26
// const DEFAULT_BOMB_COUNT = 100

func newGame(player uint, size uint, bomb uint) *Game {
	return &Game{
		Size:    size,
		Bomb:    bomb,
		Player:  make([]ClientId, player),
		Counter: Counter{},
		Turn: Turn{
			idx:  0,
			Curr: "",
		},
		Board:  GetBoard(int(size), int(bomb)),
		Winner: "",
	}
}

func (g *Game) initCounter() {
	g.Counter.BombsLeft = g.Bomb
	g.Counter.Score = make(map[ClientId]uint)
	for _, p := range g.Player {
		g.Counter.Score[p] = 0
	}
}

func (g *Game) getSize() uint {
	return g.Size
}

func (g *Game) getBombCount() uint {
	return g.Bomb
}

func (g *Game) getPlayerCap() int {
	return len(g.Player)
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
	revealables := make([]*BlockInfo, len(vertices))

	for i, v := range vertices {
		g.Board[v.Y][v.X].Visited = true
		g.Board[v.Y][v.X].VisitedBy = g.Turn.Curr
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

func (g *Game) assignTurn(cId ClientId) {
	randArr := *utils.GetRandArray(len(g.Player))
	log.Println("assignTurn Arr: ", randArr)
	for _, i := range randArr {
		if g.Player[i] == "" {
			g.Player[i] = cId
			if i == 0 {
				g.Turn.Curr = cId
			}
			return
		}
	}
}

func (g *Game) unassignTurn(cId ClientId) bool {
	for i, p := range g.Player {
		if p == cId {
			g.Player[i] = ""
			return true
		}
	}
	return false
}

func (g *Game) advanceTurn() Turn {
	g.Turn.Count++
	g.Turn.idx = (g.Turn.idx + 1) % len(g.Player)
	g.Turn.Curr = g.Player[g.Turn.idx]
	return g.Turn
}

func (g *Game) scoreCurrPlayer() (Counter, bool) {
	g.Counter.BombsLeft--
	g.Counter.Score[g.Turn.Curr]++

	isWon := g.Counter.Score[g.Turn.Curr] > g.Bomb/uint(len(g.Player))
	if isWon {
		g.Winner = g.Turn.Curr
	}
	return g.Counter, isWon
}

package ws

import (
	"encoding/json"
	"math/rand"
	"minesweeper/boardhelper"
	"time"
)

type Action struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Turn struct {
	count uint
	curr  *Client
	next  *Client
}

type GameStat struct {
	score     map[*Client]uint
	bombsLeft uint
}

type Room struct {
	id            uint
	clients       map[*Client]bool
	lobby         *Lobby
	clientCount   uint
	board         [][]boardhelper.Block
	actionHandler map[string]func(string)
	turn          Turn
	gameStat      GameStat
	update        chan *Action
	register      chan *Client
	unregister    chan *Client
}

const DEFAULT_SIZE = 26
const DEFAULT_BOMB_COUNT = 100

func newRoom(id uint, c *Client, l *Lobby) *Room {
	room := Room{
		id:          id,
		clients:     make(map[*Client]bool),
		lobby:       l,
		clientCount: 0,
		board:       boardhelper.GetBoard(DEFAULT_SIZE, DEFAULT_BOMB_COUNT),
		turn: Turn{
			count: 0,
			curr:  nil,
			next:  nil,
		},
		update:     make(chan *Action),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	// Init handler
	room.actionHandler = make(map[string]func(string))
	room.actionHandler["reveal"] = room.revealBlocks

	// Init turn-related info
	room.turn = Turn{
		count: 1,
		curr:  nil,
		next:  nil,
	}

	room.registerClient(c)
	return &room
}
func (r *Room) assignTurn(c *Client) (*Client, *Client) {
	isEmpty := r.turn.curr == nil && r.turn.next == nil
	if isEmpty {
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(2) == 0 {
			r.turn.curr = c
		} else {
			r.turn.next = c
		}
	} else if r.turn.curr == nil {
		r.turn.curr = c
	} else if r.turn.next == nil { // Not 'else'. Could be more than 3 clients in a room
		r.turn.next = c
	}
	return r.turn.curr, r.turn.next
}

func (r *Room) isPlayable(c *Client) bool {
	return r.turn.curr == c && r.turn.next != nil
}

func (r *Room) startGame() {
	type PlayerInfo struct {
		Alias string `json:"alias"`
		Score uint   `json:"score"`
	}
	type StartMsg struct {
		Id           uint       `json:"id"`
		IsPlayerTurn bool       `json:"isPlayerTurn"`
		BombsLeft    uint       `json:"bombsLeft"`
		Player       PlayerInfo `json:"player"`
		Opponent     PlayerInfo `json:"opponent"`
	}
	// Init game stat
	r.gameStat = GameStat{
		bombsLeft: DEFAULT_BOMB_COUNT,
		score:     make(map[*Client]uint),
	}
	r.gameStat.score[r.turn.curr] = 0
	r.gameStat.score[r.turn.next] = 0

	currInfo := PlayerInfo{
		Alias: r.turn.curr.alias,
		Score: 0,
	}

	nextInfo := PlayerInfo{
		Alias: r.turn.next.alias,
		Score: 0,
	}

	currStartMsg, _ := json.Marshal(StartMsg{
		Id:           r.id,
		IsPlayerTurn: true,
		BombsLeft:    DEFAULT_BOMB_COUNT,
		Player:       currInfo,
		Opponent:     nextInfo,
	})

	nextStartMsg, _ := json.Marshal(StartMsg{
		Id:           r.id,
		IsPlayerTurn: false,
		BombsLeft:    DEFAULT_BOMB_COUNT,
		Player:       nextInfo,
		Opponent:     currInfo,
	})

	r.turn.curr.update <- &Action{
		Name:    "gameStarted",
		Content: string(currStartMsg),
	}

	r.turn.next.update <- &Action{
		Name:    "gameStarted",
		Content: string(nextStartMsg),
	}
}

func (r *Room) registerClient(c *Client) {
	r.clients[c] = true
	r.assignTurn(c)
	r.clientCount++

	if r.turn.curr != nil && r.turn.next != nil {
		r.startGame()
	}
}

func (r *Room) unregisterClient(c *Client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
		if r.turn.curr == c {
			r.turn.curr = nil
		} else if r.turn.next == c {
			r.turn.next = nil
		}
	}
	if r.clientCount--; r.clientCount == 0 {
		r.lobby.unregister <- r
	}
}

func (r *Room) advanceTurn() {
	// Update turn-related info
	r.turn.count++
	r.turn.next, r.turn.curr = r.turn.curr, r.turn.next

	type TurnPassed struct {
		Count uint `json:"count"`
	}

	turn, _ := json.Marshal(TurnPassed{Count: r.turn.count})
	r.broadcast(&Action{
		Name:    "turnPassed",
		Content: string(turn),
	})
}

func (r *Room) scoreCurrPlayer() {
	// Update current player's score
	type Score struct {
		Player    uint `json:"player"`
		Opponent  uint `json:"opponent"`
		BombsLeft uint `json:"bombsLeft"`
	}
	r.gameStat.score[r.turn.curr]++
	r.gameStat.bombsLeft--

	currScore, _ := json.Marshal(Score{
		Player:    r.gameStat.score[r.turn.curr],
		Opponent:  r.gameStat.score[r.turn.next],
		BombsLeft: r.gameStat.bombsLeft,
	})

	nextScore, _ := json.Marshal(Score{
		Player:    r.gameStat.score[r.turn.next],
		Opponent:  r.gameStat.score[r.turn.curr],
		BombsLeft: r.gameStat.bombsLeft,
	})

	// Check winning condition
	if r.gameStat.score[r.turn.curr] > DEFAULT_BOMB_COUNT/2 {
		r.turn.curr.update <- &Action{
			Name:    "gameWon",
			Content: string(currScore),
		}
		r.turn.next.update <- &Action{
			Name:    "gameLost",
			Content: string(nextScore),
		}
	} else {
		r.turn.curr.update <- &Action{
			Name:    "scoreUpdated",
			Content: string(currScore),
		}
		r.turn.next.update <- &Action{
			Name:    "scoreUpdated",
			Content: string(nextScore),
		}
	}
}

func (r *Room) revealBlocks(content string) {
	// Get revealable blocks
	var v boardhelper.Vertex
	json.Unmarshal([]byte(content), &v)
	if r.board[v.Y][v.X].Visited {
		return
	}
	revealables := boardhelper.GetRevealables(&v, r.board)

	// Update visited blocks
	for _, block := range revealables {
		r.board[block.Y][block.X].Visited = true
	}
	// Advance turn or score current player
	if revealables[0].Type != boardhelper.BOMB {
		r.advanceTurn()
	} else {
		r.scoreCurrPlayer()
	}
	// Broadcast revealed blocks to clients
	data, _ := json.Marshal(struct {
		Blocks []boardhelper.BlockInfo `json:"blocks"`
	}{
		Blocks: revealables,
	})

	r.broadcast(&Action{
		Name:    "reveal",
		Content: string(data),
	})
}

func (r *Room) broadcast(action *Action) {
	for client := range r.clients {
		client.update <- action
	}
}

func (r *Room) run() {
	for {
		select {
		case c := <-r.register:
			r.registerClient(c)
		case c := <-r.unregister:
			r.unregisterClient(c)
		case action := <-r.update:
			r.actionHandler[action.Name](action.Content)
		}
	}
}

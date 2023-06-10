package ws

import (
	"encoding/json"
	"log"
	"math/rand"
	"minesweeper/boardhelper"
	"time"
)

type Turn struct {
	count uint
	curr  ClientId
	next  ClientId
}

type Counter struct {
	score     map[ClientId]uint
	bombsLeft uint
}

type PlayerInfo struct {
	Id     ClientId `json:"id"`
	Alias  string   `json:"alias"`
	Score  uint     `json:"score"`
	IsTurn bool     `json:"isTurn"`
}

type GameStat struct {
	BombsLeft uint                    `json:"bombsLeft"`
	Players   map[ClientId]PlayerInfo `json:"players"`
	Visible   []boardhelper.BlockInfo `json:"visible"`
}

type Room struct {
	id            uint
	clients       map[ClientId]*Client
	lobby         *Lobby
	board         [][]boardhelper.Block
	actionHandler map[string]func(string)
	turn          Turn
	counter       Counter
	update        chan *Action
	register      chan *Client
	unregister    chan *Client
	disconnect    chan *Client
	reconnect     chan *Client
	timeouts      map[ClientId]int64
	stop          chan bool
}

const DEFAULT_SIZE = 26
const DEFAULT_BOMB_COUNT = 100
const TIMELIMIT_IN_SEC = 5 // 5 minutes

func newRoom(id uint, c *Client, l *Lobby) *Room {
	r := &Room{
		id:      id,
		clients: make(map[ClientId]*Client),
		lobby:   l,
		board:   boardhelper.GetBoard(DEFAULT_SIZE, DEFAULT_BOMB_COUNT),
		turn: Turn{
			count: 1,
			curr:  "",
			next:  "",
		},
		update:     make(chan *Action),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		disconnect: make(chan *Client),
		reconnect:  make(chan *Client),
		timeouts:   make(map[ClientId]int64),
	}
	// Init handler
	r.actionHandler = make(map[string]func(string))
	r.actionHandler["reveal"] = r.revealBlocks
	go r.run()

	// Register Client
	r.registerClient(c)

	// Update client
	c.update <- &Action{
		Name:    "roomCreated",
		Content: "{}",
	}
	return r
}

func (r *Room) registerClient(c *Client) {
	r.clients[c.id] = c
	curr, next := r.assignTurn(c)

	type RoomIdMsg struct {
		Id uint `json:"id"`
	}
	rId, _ := json.Marshal(RoomIdMsg{Id: r.id})

	c.update <- &Action{
		Name:    "roomId",
		Content: string(rId),
	}

	if r.clients[curr] == nil || r.clients[next] == nil {
		return
	}
	if curr == c.id || next == c.id {
		r.startGame()
	} else {
		// FIXME: Add watch game for non-player
		// r.watchGame()
	}
}

func (r *Room) unregisterClient(c *Client) {
	if _, ok := r.clients[c.id]; ok {
		delete(r.clients, c.id)
		if r.clients[r.turn.curr] == c {
			r.turn.curr = ""
		} else if r.clients[r.turn.next] == c {
			r.turn.next = ""
		}
	}
	if len(r.clients) == 0 {
		r.stop <- true
		r.lobby.unregister <- r
	}
}

func (r *Room) disconnectClient(c *Client) {
	log.Println("Client", c.alias, "disconnected")
	if _, ok := r.timeouts[c.id]; !ok {
		r.timeouts[c.id] = time.Now().Unix()
	}
}

func (r *Room) reconnectClient(c *Client) {
	// Reconnect user if id is in timeout map
	_, timeoutOk := r.timeouts[c.id]
	_, clientsOk := r.clients[c.id]
	if !timeoutOk || !clientsOk {
		log.Println("Cannot reconnect client")
		c.update <- &Action{
			Name:    "reconnFailed",
			Content: "{}",
		}
		return
	}
	// Remove from timeout map
	delete(r.timeouts, c.id)

	// Reuse alias before disconnection, reassign
	c.alias = r.clients[c.id].alias
	r.clients[c.id] = c
	c.room = r

	log.Println("Client", c.alias, "reconnected")

	// Return Game Stat
	stat, _ := json.Marshal(r.getGameStat())

	c.update <- &Action{
		Name:    "gameStat",
		Content: string(stat),
	}
}

func (r *Room) assignTurn(c *Client) (ClientId, ClientId) {
	isEmpty := r.turn.curr == "" && r.turn.next == ""
	if isEmpty {
		if rand.Intn(2) == 0 {
			r.turn.curr = c.id
		} else {
			r.turn.next = c.id
		}
	} else if r.turn.curr == "" {
		r.turn.curr = c.id
	} else if r.turn.next == "" { // Not 'else'. Could be more than 3 clients in a room
		r.turn.next = c.id
	}
	return r.turn.curr, r.turn.next
}

func (r *Room) isPlayable(c *Client) bool {
	return r.turn.curr == c.id && r.turn.next != ""
}

func (r *Room) startGame() {

	// Init game stat
	r.counter = Counter{
		bombsLeft: DEFAULT_BOMB_COUNT,
		score:     make(map[ClientId]uint),
	}
	r.counter.score[r.turn.curr] = 0
	r.counter.score[r.turn.next] = 0

	gameStatMsg, _ := json.Marshal(r.getGameStat())

	r.broadcast(&Action{
		Name:    "gameStat",
		Content: string(gameStatMsg),
	})
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
	r.counter.score[r.turn.curr]++
	r.counter.bombsLeft--

	currScore, _ := json.Marshal(Score{
		Player:    r.counter.score[r.turn.curr],
		Opponent:  r.counter.score[r.turn.next],
		BombsLeft: r.counter.bombsLeft,
	})

	nextScore, _ := json.Marshal(Score{
		Player:    r.counter.score[r.turn.next],
		Opponent:  r.counter.score[r.turn.curr],
		BombsLeft: r.counter.bombsLeft,
	})

	// Check winning condition
	if r.counter.score[r.turn.curr] > DEFAULT_BOMB_COUNT/2 {
		r.clients[r.turn.curr].update <- &Action{
			Name:    "gameWon",
			Content: string(currScore),
		}
		r.clients[r.turn.next].update <- &Action{
			Name:    "gameLost",
			Content: string(nextScore),
		}
	} else {
		r.clients[r.turn.curr].update <- &Action{
			Name:    "scoreUpdated",
			Content: string(currScore),
		}
		r.clients[r.turn.next].update <- &Action{
			Name:    "scoreUpdated",
			Content: string(nextScore),
		}
	}
}

func (r *Room) getGameStat() *GameStat {
	curr, next := r.turn.curr, r.turn.next
	if curr == "" || next == "" {
		return &GameStat{}
	}
	currInfo := PlayerInfo{
		Id:     r.clients[r.turn.curr].id,
		Alias:  r.clients[r.turn.curr].alias,
		Score:  r.counter.score[r.turn.curr],
		IsTurn: true,
	}

	nextInfo := PlayerInfo{
		Id:     r.clients[next].id,
		Alias:  r.clients[next].alias,
		Score:  r.counter.score[next],
		IsTurn: false,
	}

	return &GameStat{
		BombsLeft: r.counter.bombsLeft,
		Players:   map[ClientId]PlayerInfo{currInfo.Id: currInfo, nextInfo.Id: nextInfo},
		Visible:   r.getVisibleBlocks(),
	}
}

func (r *Room) getVisibleBlocks() []boardhelper.BlockInfo {
	s := []boardhelper.BlockInfo{}
	for i := range r.board {
		for j := range r.board {
			if r.board[i][j].Visited {
				s = append(s, boardhelper.BlockInfo{
					X:     j,
					Y:     i,
					Block: r.board[i][j],
				})
			}
		}
	}
	return s
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

func (r *Room) checkActivity(now int64) {
	for cId, t := range r.timeouts {
		if now-t > TIMELIMIT_IN_SEC {
			log.Println("Removing client", r.clients[cId].alias, "from room", r.id)
			delete(r.clients, cId)
			delete(r.timeouts, cId)
		}
	}
}

func (r *Room) broadcast(action *Action) {
	for client := range r.clients {
		r.clients[client].update <- action
	}
}

func (r *Room) run() {
	// Setup ticker
	ticker := time.NewTicker((time.Second))
	doneChecking := make(chan bool)
	defer func() {
		doneChecking <- true
		ticker.Stop()
	}()

	// Check for client activity
	go func() {
		for {
			select {
			case t := <-ticker.C:
				r.checkActivity(t.Unix())
			case <-doneChecking:
				break
			}
		}
	}()

	for {
		select {
		case action := <-r.update:
			r.actionHandler[action.Name](action.Content)
		case c := <-r.register:
			r.registerClient(c)
		case c := <-r.unregister:
			r.unregisterClient(c)
		case c := <-r.disconnect:
			r.disconnectClient(c)
		case c := <-r.reconnect:
			r.reconnectClient(c)
		case <-r.stop:
			break
		}
	}
}

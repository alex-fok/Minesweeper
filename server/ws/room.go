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

type GameStat struct {
	score     map[ClientId]uint
	bombsLeft uint
}

type Room struct {
	id            uint
	clients       map[ClientId]*Client
	lobby         *Lobby
	board         [][]boardhelper.Block
	actionHandler map[string]func(string)
	turn          Turn
	gameStat      GameStat
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

	// Update client with room id
	content, _ := json.Marshal(struct {
		RoomId uint `json:"roomId"`
	}{
		RoomId: id,
	})
	c.update <- &Action{
		Name:    "roomCreated",
		Content: string(content),
	}
	return r
}

func (r *Room) registerClient(c *Client) {
	r.clients[c.id] = c
	r.assignTurn(c)

	if r.clients[r.turn.curr] != nil && r.clients[r.turn.next] != nil {
		r.startGame()
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
			Name:    "Failed Reconnection",
			Content: "{}",
		}
		return
	}
	// Remove from timeout map
	delete(r.timeouts, c.id)

	// Reuse alias before disconnection, reassign
	c.alias = r.clients[c.id].alias
	r.clients[c.id] = c

	log.Println("Client", c.alias, "reconnected")
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
		score:     make(map[ClientId]uint),
	}
	r.gameStat.score[r.turn.curr] = 0
	r.gameStat.score[r.turn.next] = 0

	currInfo := PlayerInfo{
		Alias: r.clients[r.turn.curr].alias,
		Score: 0,
	}

	nextInfo := PlayerInfo{
		Alias: r.clients[r.turn.next].alias,
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

	r.clients[r.turn.curr].update <- &Action{
		Name:    "gameStarted",
		Content: string(currStartMsg),
	}

	r.clients[r.turn.next].update <- &Action{
		Name:    "gameStarted",
		Content: string(nextStartMsg),
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

package game

import (
	"encoding/json"
	"minesweeper/types"
	"sync"
	"time"
)

type Action = types.Action
type Client = types.ClientMeta
type Content = string

type BoardConfig struct {
	Size uint `json:"size"`
	Bomb uint `json:"bomb"`
}

type PlayerInfo struct {
	Id       ClientId `json:"id"`
	Alias    string   `json:"alias"`
	Score    uint     `json:"score"`
	IsTurn   bool     `json:"isTurn"`
	IsOnline bool     `json:"isOnline"`
}

type OnlineStatus = types.OnlineStatus

type GameStat struct {
	BoardConfig *BoardConfig             `json:"boardConfig"`
	BombsLeft   uint                     `json:"bombsLeft"`
	Players     map[ClientId]*PlayerInfo `json:"players"`
	Visible     []BlockInfo              `json:"visible"`
}

type Timer struct {
	limit float64
	reset chan bool
	stop  chan bool
}

type Driver struct {
	mu         sync.Mutex
	game       *Game
	timer      Timer
	stopTimer  chan bool
	playerCap  int
	players    map[ClientId]*Client
	rematchReq map[ClientId]bool
	broadcast  func(action *Action)
	lastPlayed time.Time
	isDone     bool
}

func NewDriver(timeLimit uint, size uint, bomb uint) *Driver {
	d := Driver{
		game: newGame(size, bomb),
		timer: Timer{
			limit: float64(timeLimit),
			reset: make(chan bool),
			stop:  make(chan bool),
		},
		playerCap:  2,
		players:    make(map[ClientId]*Client),
		rematchReq: make(map[ClientId]bool),
		isDone:     false,
	}
	return &d
}

func (d *Driver) RegisterPlayer(c *Client) bool {
	isPlayer := len(d.players) < d.playerCap

	if isPlayer {
		d.game.assignTurn(c.Id)
		d.players[c.Id] = c
	}
	return isPlayer
}

func (d *Driver) UnregisterPlayer(cId ClientId) {
	currId, nextId := d.game.unassignTurn(cId)
	if currId == "" || nextId == "" {
		d.broadcast(&Action{
			Name:    "gameEnded",
			Content: "{}",
		})
	}
}

func (d *Driver) DisconnectPlayer(cId ClientId) bool {
	if _, ok := d.players[cId]; ok {
		if d.timer.limit != 0 && d.IsGameReady() {
			d.stopTimer <- true
		}
		d.players[cId].IsOnline = false
		d.players[cId].IsReady = false
		type PlayerOnline struct {
			Player   ClientId `json:"player"`
			IsOnline bool     `json:"isOnline"`
		}
		disconnPlayer, _ := json.Marshal(PlayerOnline{
			Player:   cId,
			IsOnline: false,
		})
		d.broadcast(&Action{
			Name:    "playerOnline",
			Content: string(disconnPlayer),
		})
		return true
	}
	return false
}

func (d *Driver) ReconnectPlayer(cId ClientId) {
	if _, ok := d.players[cId]; ok {
		d.players[cId].IsOnline = true
		type PlayerOnline struct {
			Player   ClientId `json:"player"`
			IsOnline bool     `json:"isOnline"`
		}
		reconnPlayer, _ := json.Marshal(PlayerOnline{
			Player:   cId,
			IsOnline: true,
		})
		d.broadcast(&Action{
			Name:    "playerOnline",
			Content: string(reconnPlayer),
		})
	}
}

func (d *Driver) GetGameStat() *GameStat {
	counter, turn := d.game.getCounter(), d.game.getTurn()

	gameStat := GameStat{
		BoardConfig: &BoardConfig{
			Size: d.game.getSize(),
			Bomb: d.game.getBombCount(),
		},
		BombsLeft: counter.BombsLeft,
		Players:   make(map[ClientId]*PlayerInfo),
		Visible:   d.game.getVisibleBlocks(),
	}
	for _, p := range d.players {
		gameStat.Players[p.Id] = &PlayerInfo{
			Id:       p.Id,
			Alias:    p.Alias,
			Score:    counter.Score[p.Id],
			IsTurn:   turn.Curr == p.Id,
			IsOnline: p.IsOnline,
		}
	}
	return &gameStat
}

func (d *Driver) GetPlayerOnlineStatus() map[ClientId]OnlineStatus {
	result := make(map[ClientId]OnlineStatus)
	for k, v := range d.players {
		result[k] = OnlineStatus{
			Alias:    v.Alias,
			IsOnline: v.IsOnline,
			IsReady:  v.IsReady,
		}
	}
	return result
}

func (d *Driver) GetPlayerCap() int {
	return d.playerCap
}

func (d *Driver) SetBroadcast(b func(*Action)) {
	d.broadcast = b
}

func (d *Driver) IsGameReady() bool {
	if len(d.players) != d.playerCap {
		return false
	}
	for _, p := range d.players {
		if !p.IsReady {
			return false
		}
	}
	return true
}

func (d *Driver) isExpired() bool {
	return d.timer.limit != 0 && time.Since(d.lastPlayed).Seconds() > (d.timer.limit-.1)
}

func (d *Driver) updateLastPlayed() {
	d.lastPlayed = time.Now()
}

func (d *Driver) startTimer() {
	d.stopTimer = make(chan bool)
	timesup := make(chan bool)
	timer := func() {
		time.Sleep(time.Second * time.Duration(int64(d.timer.limit)))
		if d.isExpired() {
			timesup <- true
		}
	}
	go timer()
	for {
		select {
		case <-d.stopTimer:
			return
		case <-timesup:
			d.mu.Lock()
			if d.isExpired() {
				d.advanceTurn()
				d.updateLastPlayed()
				go timer()
			}
			d.mu.Unlock()
		case <-d.timer.reset:
			go timer()
		}

	}
}

func (d *Driver) advanceTurn() {
	turn := d.game.advanceTurn()
	type TurnPassed struct {
		Count uint     `json:"count"`
		Curr  ClientId `json:"curr"`
	}

	turnPassed, _ := json.Marshal(TurnPassed{
		Count: turn.Count,
		Curr:  turn.Curr,
	})
	d.broadcast(&Action{
		Name:    "turnPassed",
		Content: string(turnPassed),
	})
}

func (d *Driver) scoreCurrPlayer() {
	counter, isWon := d.game.scoreCurrPlayer()

	scoreUpdated, _ := json.Marshal(counter)
	d.broadcast(&Action{
		Name:    "scoreUpdated",
		Content: string(scoreUpdated),
	})

	if isWon {
		type GameEnded struct {
			Winner ClientId `json:"winner"`
		}

		if d.timer.limit != 0 {
			d.stopTimer <- true
		}
		gameEnded, _ := json.Marshal(GameEnded{Winner: d.game.getWinner()})
		d.broadcast(&Action{
			Name:    "gameEnded",
			Content: string(gameEnded),
		})
	}
}

func (d *Driver) reveal(cId ClientId, content string) {
	defer func() {
		d.mu.Unlock()
	}()

	// Lock
	d.mu.Lock()

	isDone := d.isDone
	expired := d.isExpired()
	notTurn := cId != d.game.getTurn().Curr

	if isDone || expired || notTurn {
		return
	}

	// Get revealable blocks
	var v Vertex
	json.Unmarshal([]byte(content), &v)

	if d.game.Board[v.Y][v.X].Visited {
		return
	}
	revealables := d.game.reveal(&v)

	type Revealed struct {
		Blocks []*BlockInfo `json:"blocks"`
	}

	data, _ := json.Marshal(Revealed{
		Blocks: revealables,
	})
	d.broadcast(&Action{
		Name:    "reveal",
		Content: string(data),
	})
	// Advance turn or score current player
	if revealables[0].Type != BOMB {
		d.advanceTurn()
	} else {
		d.scoreCurrPlayer()
	}

	// If time limit is set, reset timer

	if d.timer.limit != 0 {
		d.updateLastPlayed()
		d.timer.reset <- true
	}
}

func (d *Driver) rematch(cId ClientId, content string) {
	if _, ok := d.players[cId]; d.isDone || !ok {
		return
	}

	type Req struct {
		Rematch bool `json:"rematch"`
	}

	var r Req
	json.Unmarshal([]byte(content), &r)
	d.rematchReq[cId] = r.Rematch

	// If player decides not to rematch, game is closed
	if !d.rematchReq[cId] {
		d.isDone = true
		type Message struct {
			Message string `json:"message"`
		}
		msg, _ := json.Marshal(Message{
			Message: "Player " + d.players[cId].Alias + " declined the rematch.",
		})
		d.broadcast(&Action{
			Name:    "message",
			Content: string(msg),
		})
	}
	// Check if all players select "Rematch"
	rematch := true
	for _, client := range d.players {
		if _, ok := d.rematchReq[client.Id]; !ok {
			rematch = false
			break
		}
	}

	if rematch {
		d.game = newGame(d.game.getSize(), d.game.getBombCount())
		for _, player := range d.players {
			d.game.assignTurn(player.Id)
		}
		d.game.initCounter()
		d.rematchReq = make(map[ClientId]bool)
		gameStatMsg, _ := json.Marshal(d.GetGameStat())

		d.broadcast(&Action{
			Name:    "gameStat",
			Content: string(gameStatMsg),
		})
	}
}

func (d *Driver) playerReady(cId ClientId, content string) {
	type Req struct {
		IsReady bool `json:"isReady"`
	}
	var r Req
	json.Unmarshal([]byte(content), &r)

	if _, ok := d.players[cId]; ok {
		d.players[cId].IsReady = r.IsReady

		type PlayerReady struct {
			Player  ClientId `json:"player"`
			IsReady bool     `json:"isReady"`
		}

		playerReady, _ := json.Marshal(&PlayerReady{
			Player:  cId,
			IsReady: r.IsReady,
		})
		d.broadcast(&Action{
			Name:    "playerReady",
			Content: string(playerReady),
		})
		if d.IsGameReady() {
			d.StartGame()
		}
	}
}

func (d *Driver) RenamePlayer(cId ClientId, alias string) {
	if player, ok := d.players[cId]; !d.isDone && ok {
		player.Alias = alias
	}
}

func (d *Driver) HandleGameUpdate(cId ClientId, update *Action) {
	switch update.Name {
	case "reveal":
		d.reveal(cId, update.Content)
	case "rematch":
		d.rematch(cId, update.Content)
	case "ready":
		d.playerReady(cId, update.Content)
	}
}

func (d *Driver) StartGame() {
	// Init game stat
	d.game.initCounter()
	gameStat, _ := json.Marshal(d.GetGameStat())
	d.updateLastPlayed()

	if d.timer.limit != 0 {
		go d.startTimer()
	}

	d.broadcast(&Action{
		Name:    "gameStat",
		Content: string(gameStat),
	})
}

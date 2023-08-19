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
	TimeLimit   float64                  `json:"timeLimit"`
	LastPlayed  time.Time                `json:"lastPlayed"`
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
	players    map[ClientId]*Client
	rematchReq map[ClientId]bool
	broadcast  func(action *Action)
	lastPlayed time.Time
	isStarted  bool
	isDone     bool
}

func NewDriver(timeLimit uint, player uint, size uint, bomb uint) *Driver {
	d := Driver{
		game: newGame(player, size, bomb),
		timer: Timer{
			limit: float64(timeLimit),
			reset: make(chan bool),
			stop:  make(chan bool),
		},
		players:    make(map[ClientId]*Client),
		rematchReq: make(map[ClientId]bool),
		isDone:     false,
	}
	return &d
}

func (d *Driver) RegisterPlayer(c *Client) bool {
	isPlayer := len(d.players) < d.game.getPlayerCap()

	if isPlayer {
		d.game.assignTurn(c.Id)
		d.players[c.Id] = c
	}
	return isPlayer
}

func (d *Driver) UnregisterPlayer(cId ClientId) {
	isUnassigned := d.game.unassignTurn(cId)
	type GameEnded struct {
		IsCanceled bool     `json:"isCanceled"`
		Winner     ClientId `json:"winner"`
	}
	gameEnded, _ := json.Marshal(&GameEnded{
		IsCanceled: true,
		Winner:     "",
	})
	if isUnassigned {
		d.broadcast(&Action{
			Name:    "gameEnded",
			Content: string(gameEnded),
		})
	}
}

func (d *Driver) DisconnectPlayer(cId ClientId) bool {
	if _, ok := d.players[cId]; ok {
		if d.timer.limit != 0 && d.IsGameReady() {
			d.timer.stop <- true
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
		BombsLeft:  counter.BombsLeft,
		Players:    make(map[ClientId]*PlayerInfo),
		Visible:    d.game.getVisibleBlocks(),
		TimeLimit:  d.timer.limit,
		LastPlayed: d.lastPlayed,
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
	return d.game.getPlayerCap()
}

func (d *Driver) SetBroadcast(b func(*Action)) {
	d.broadcast = b
}

func (d *Driver) IsGameReady() bool {
	if len(d.players) != d.game.getPlayerCap() {
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

func (d *Driver) updateLastPlayed() time.Time {
	d.lastPlayed = time.Now()
	return d.lastPlayed
}

func (d *Driver) startTimer() {
	var timer *time.Timer

	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()
	resetTimer := func() {
		if timer != nil {
			timer.Stop()
		}
		timer = time.NewTimer(time.Second * time.Duration(int64(d.timer.limit)))
	}
	timer = time.NewTimer(time.Second * time.Duration(int64(d.timer.limit)))
	for {
		select {
		case <-d.timer.stop:
			return
		case <-timer.C:
			d.mu.Lock()
			if d.isExpired() {
				d.advanceTurn()
				resetTimer()
			}
			d.mu.Unlock()
		case <-d.timer.reset:
			resetTimer()
		}
	}
}

func (d *Driver) advanceTurn() {
	turn := d.game.advanceTurn()
	type TurnPassed struct {
		Count      uint      `json:"count"`
		Curr       ClientId  `json:"curr"`
		LastPlayed time.Time `json:"lastPlayed"`
	}
	turnPassed, _ := json.Marshal(TurnPassed{
		Count:      turn.Count,
		Curr:       turn.Curr,
		LastPlayed: d.updateLastPlayed(),
	})
	d.broadcast(&Action{
		Name:    "turnPassed",
		Content: string(turnPassed),
	})
}

func (d *Driver) scoreCurrPlayer() {
	counter, isWon := d.game.scoreCurrPlayer()
	type Score struct {
		BombsLeft  uint              `json:"bombsLeft"`
		Score      map[ClientId]uint `json:"score"`
		LastPlayed time.Time         `json:"lastPlayed"`
	}
	counterWithDate := Score{
		BombsLeft:  counter.BombsLeft,
		Score:      counter.Score,
		LastPlayed: d.updateLastPlayed(),
	}
	scoreUpdated, _ := json.Marshal(counterWithDate)
	d.broadcast(&Action{
		Name:    "scoreUpdated",
		Content: string(scoreUpdated),
	})

	if isWon || counter.BombsLeft == 0 {
		if d.timer.limit != 0 {
			d.timer.stop <- true
		}
		type GameEnded struct {
			IsCanceled bool     `json:"isCanceled"`
			Winner     ClientId `json:"winner"`
		}
		gameEnded, _ := json.Marshal(GameEnded{
			IsCanceled: false,
			Winner:     d.game.getWinner(),
		})
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
	isExpired := d.isExpired()
	notTurn := cId != d.game.getTurn().Curr

	if isDone || isExpired || notTurn {
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
		d.game = newGame(d.game.getSize(), d.game.getBombCount(), uint(d.game.getPlayerCap()))
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
	if !d.isStarted {
		d.isStarted = true
		// Init game stat
		d.game.initCounter()
	}

	d.updateLastPlayed()
	gameStat, _ := json.Marshal(d.GetGameStat())

	if d.timer.limit != 0 {
		go d.startTimer()
	}

	d.broadcast(&Action{
		Name:    "gameStat",
		Content: string(gameStat),
	})
}

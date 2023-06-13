package game

import (
	"encoding/json"
	"minesweeper/types"
)

type Action = types.Action
type Client = types.ClientMeta

type PlayerInfo struct {
	Id       ClientId `json:"id"`
	Alias    string   `json:"alias"`
	Score    uint     `json:"score"`
	IsTurn   bool     `json:"isTurn"`
	IsOnline bool     `json:"isOnline"`
}

type GameStat struct {
	BombsLeft uint                     `json:"bombsLeft"`
	Players   map[ClientId]*PlayerInfo `json:"players"`
	Visible   []BlockInfo              `json:"visible"`
}

type Driver struct {
	ActionHandler map[string]func(ClientId, string) []*Action
	game          *Game
	Players       map[ClientId]*Client
	RematchReq    map[ClientId]bool
	isDone        bool
}

func NewDriver() *Driver {
	d := Driver{
		ActionHandler: make(map[string]func(ClientId, string) []*Action),
		game:          newGame(),
		Players:       make(map[ClientId]*Client),
		RematchReq:    make(map[ClientId]bool),
		isDone:        false,
	}
	d.ActionHandler["reveal"] = d.Reveal
	return &d
}

func (d *Driver) RegisterPlayer(c *Client) []*Action {
	actions := []*Action{}
	currId, nextId := d.game.assignTurn(c.Id)

	isGameReady := currId != "" && nextId != ""
	isPlayer := currId == c.Id || nextId == c.Id

	if isPlayer {
		d.Players[c.Id] = c
		if !isGameReady {
			return actions
		}
		actions = append(actions, d.StartGame())
	} else {
		if !isGameReady {
			return actions
		}
		gameStatMsg, _ := json.Marshal(d.GetGameStat())

		actions = append(actions, &Action{
			Name:    "gameStat",
			Content: string(gameStatMsg),
		})
	}
	return actions
}

func (d *Driver) UnregisterPlayer(cId ClientId) []*Action {
	actions := []*Action{}
	currId, nextId := d.game.unassignTurn(cId)
	if currId == "" || nextId == "" {
		actions = append(actions, &Action{
			Name:    "gameEnded",
			Content: "{}",
		})
	}
	return actions
}

func (d *Driver) DisconnectPlayer(cId ClientId) *Action {
	var action *Action = nil
	if _, ok := d.Players[cId]; ok {
		d.Players[cId].IsOnline = false
		type DisconnPlayer struct {
			Client ClientId `json:"client"`
		}
		disconnPlayer, _ := json.Marshal(DisconnPlayer{
			Client: cId,
		})
		action = &Action{
			Name:    "playerOffline",
			Content: string(disconnPlayer),
		}
	}
	return action
}

func (d *Driver) ReconnectPlayer(cId ClientId) *Action {
	var action *Action = nil
	if _, ok := d.Players[cId]; ok {
		d.Players[cId].IsOnline = true
		type ReconnPlayer struct {
			Client ClientId `json:"client"`
		}
		reconnPlayer, _ := json.Marshal(ReconnPlayer{
			Client: cId,
		})
		action = &Action{
			Name:    "playerOnline",
			Content: string(reconnPlayer),
		}
	}
	return action
}

func (d *Driver) advanceTurn() []*Action {
	actions := []*Action{}
	turn := d.game.advanceTurn()
	type TurnPassed struct {
		Count uint     `json:"count"`
		Curr  ClientId `json:"curr"`
	}

	turnPassed, _ := json.Marshal(TurnPassed{
		Count: turn.Count,
		Curr:  turn.Curr,
	})
	actions = append(actions, &Action{
		Name:    "turnPassed",
		Content: string(turnPassed),
	})
	return actions
}

func (d *Driver) scoreCurrPlayer() []*Action {
	actions := []*Action{}
	counter, isWon := d.game.scoreCurrPlayer()

	scoreUpdated, _ := json.Marshal(counter)
	actions = append(actions, &Action{
		Name:    "scoreUpdated",
		Content: string(scoreUpdated),
	})

	if isWon {
		type GameEnded struct {
			Winner ClientId `json:"winner"`
		}
		gameEnded, _ := json.Marshal(GameEnded{Winner: d.game.getWinner()})
		actions = append(actions, &Action{
			Name:    "gameEnded",
			Content: string(gameEnded),
		})
	}
	return actions
}

func (d *Driver) Reveal(cId ClientId, content string) []*Action {
	actions := []*Action{}

	if d.isDone || cId != d.game.getTurn().Curr {
		return actions
	}

	// Get revealable blocks
	var v Vertex
	json.Unmarshal([]byte(content), &v)

	if d.game.Board[v.Y][v.X].Visited {
		return actions
	}
	revealables := d.game.reveal(&v)

	type Revealed struct {
		Blocks []*BlockInfo `json:"blocks"`
	}

	data, _ := json.Marshal(Revealed{
		Blocks: revealables,
	})
	actions = append(actions, &Action{
		Name:    "reveal",
		Content: string(data),
	})
	// Advance turn or score current player
	var a []*Action
	if revealables[0].Type != BOMB {
		a = d.advanceTurn()
	} else {
		a = d.scoreCurrPlayer()
	}
	actions = append(actions, a...)
	return actions
}

func (d *Driver) Rematch(cId ClientId, content string) []*Action {
	actions := []*Action{}

	if _, ok := d.Players[cId]; d.isDone || !ok {
		return actions
	}

	type Req struct {
		Rematch bool `json:"rematch"`
	}
	var r Req
	json.Unmarshal([]byte(content), &r)
	d.RematchReq[cId] = r.Rematch

	// If player decides not to rematch, game is closed
	if !d.RematchReq[cId] {
		d.isDone = true
		type Message struct {
			Message string `json:"message"`
		}
		msg, _ := json.Marshal(Message{
			Message: "Game is closed",
		})
		actions = append(actions, &Action{
			Name:    "message",
			Content: string(msg),
		})
		return actions
	}
	// Check if all players select "Rematch"
	rematch := true
	for _, client := range d.Players {
		if _, ok := d.RematchReq[client.Id]; !ok {
			rematch = false
			break
		}
	}

	if rematch {
		d.game = newGame()
		for _, player := range d.Players {
			d.game.assignTurn(player.Id)
		}
		d.game.initCounter()
		d.RematchReq = make(map[ClientId]bool)
		gameStatMsg, _ := json.Marshal(d.GetGameStat())

		actions = append(actions, &Action{
			Name:    "gameStat",
			Content: string(gameStatMsg),
		})
	}
	return actions
}

func (d *Driver) GetGameStat() *GameStat {
	counter, turn := d.game.getCounter(), d.game.getTurn()

	gameStat := GameStat{
		BombsLeft: counter.BombsLeft,
		Players:   make(map[ClientId]*PlayerInfo),
		Visible:   d.game.getVisibleBlocks(),
	}
	for _, p := range d.Players {
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

func (d *Driver) StartGame() *Action {

	// Init game stat
	d.game.initCounter()

	gameStatMsg, _ := json.Marshal(d.GetGameStat())

	return &Action{
		Name:    "gameStat",
		Content: string(gameStatMsg),
	}
}

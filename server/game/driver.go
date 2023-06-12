package game

import (
	"encoding/json"
	"minesweeper/types"
)

type Action = types.Action
type Client = types.ClientMeta

type PlayerInfo struct {
	Id     ClientId `json:"id"`
	Alias  string   `json:"alias"`
	Score  uint     `json:"score"`
	IsTurn bool     `json:"isTurn"`
}

type GameStat struct {
	BombsLeft uint                     `json:"bombsLeft"`
	Players   map[ClientId]*PlayerInfo `json:"players"`
	Visible   []BlockInfo              `json:"visible"`
}

type Driver struct {
	ActionHandler map[string]func(ClientId, string) []*Action
	Game          Game
	Players       map[ClientId]*Client
}

func NewDriver() *Driver {
	d := Driver{
		ActionHandler: make(map[string]func(ClientId, string) []*Action),
		Game:          *newGame(),
		Players:       make(map[ClientId]*Client),
	}
	d.ActionHandler["reveal"] = d.Reveal
	return &d
}

func (d *Driver) RegisterPlayer(c *Client) []*Action {
	actions := []*Action{}
	currId, nextId := d.Game.assignTurn(c.Id)

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
	currId, nextId := d.Game.unassignTurn(cId)
	if currId == "" || nextId == "" {
		actions = append(actions, &Action{
			Name:    "gameEnded",
			Content: "{}",
		})
	}
	return actions
}

func (d *Driver) advanceTurn() []*Action {
	actions := []*Action{}
	turn := d.Game.advanceTurn()
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
	counter, isWon := d.Game.scoreCurrPlayer()

	scoreUpdated, _ := json.Marshal(counter)
	actions = append(actions, &Action{
		Name:    "scoreUpdated",
		Content: string(scoreUpdated),
	})

	if isWon {
		type GameEnded struct {
			Winner ClientId `json:"winner"`
		}
		gameEnded, _ := json.Marshal(GameEnded{Winner: d.Game.getWinner()})
		actions = append(actions, &Action{
			Name:    "gameEnded",
			Content: string(gameEnded),
		})
	}
	return actions
}

func (d *Driver) Reveal(cId ClientId, content string) []*Action {
	actions := []*Action{}

	if cId != d.Game.getTurn().Curr {
		return actions
	}

	// Get revealable blocks
	var v Vertex
	json.Unmarshal([]byte(content), &v)

	if d.Game.Board[v.Y][v.X].Visited {
		return actions
	}
	revealables := GetRevealables(&v, d.Game.getBoard())

	// Update visited blocks
	for _, block := range revealables {
		d.Game.Board[block.Y][block.X].Visited = true
	}

	type Revealed struct {
		Blocks []BlockInfo `json:"blocks"`
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

func (d *Driver) GetGameStat() *GameStat {
	counter, turn := d.Game.getCounter(), d.Game.getTurn()

	gameStat := GameStat{
		BombsLeft: counter.BombsLeft,
		Players:   make(map[ClientId]*PlayerInfo),
		Visible:   d.Game.getVisibleBlocks(),
	}
	for _, p := range d.Players {
		gameStat.Players[p.Id] = &PlayerInfo{
			Id:     p.Id,
			Alias:  p.Alias,
			Score:  counter.Score[p.Id],
			IsTurn: turn.Curr == p.Id,
		}
	}
	return &gameStat
}

func (d *Driver) StartGame() *Action {

	// Init game stat
	d.Game.initCounter()

	gameStatMsg, _ := json.Marshal(d.GetGameStat())

	return &Action{
		Name:    "gameStat",
		Content: string(gameStatMsg),
	}
}

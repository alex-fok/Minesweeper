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
	ActionHandler map[string]func(string) []*Action
	Game          Game
	Players       map[ClientId]*Client
}

func NewDriver() *Driver {
	d := Driver{
		ActionHandler: make(map[string]func(string) []*Action),
		Game:          *CreateGame(),
		Players:       make(map[ClientId]*Client),
	}
	d.ActionHandler["reveal"] = d.Reveal
	return &d
}

func (d *Driver) RegisterPlayer(c *Client) []*Action {
	actions := []*Action{}
	currId, nextId := d.Game.AssignTurn(c.Id)

	isGameReady := currId != "" && nextId != ""
	isPlayer := currId == c.Id || nextId == c.Id

	if isPlayer {
		d.Players[c.Id] = c
		if !isGameReady {
			return actions
		}
		actions = append(actions, d.startGame())
	} else {
		if !isGameReady {
			return actions
		}
		// Give viewer game stat
		// FIXME: Add watch game function for non-player
		// r.watchGame()
	}
	return actions
}

func (d *Driver) UnregisterPlayer(cId ClientId) []*Action {
	actions := []*Action{}
	currId, nextId := d.Game.UnassignTurn(cId)
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
	turn := d.Game.AdvanceTurn()
	type TurnPassed struct {
		Count uint `json:"count"`
	}

	turnPassed, _ := json.Marshal(TurnPassed{Count: turn.Count})
	actions = append(actions, &Action{
		Name:    "turnPassed",
		Content: string(turnPassed),
	})
	return actions
}

func (d *Driver) scoreCurrPlayer() []*Action {
	actions := []*Action{}
	counter, isWon := d.Game.ScoreCurrPlayer()

	scoreUpdated, _ := json.Marshal(counter)
	actions = append(actions, &Action{
		Name:    "scoreUpdated",
		Content: string(scoreUpdated),
	})

	if isWon {
		type GameEnded struct {
			Winner ClientId `json:"winner"`
		}
		gameEnded, _ := json.Marshal(GameEnded{Winner: d.Game.GetWinner()})
		actions = append(actions, &Action{
			Name:    "gameEnded",
			Content: string(gameEnded),
		})
	}
	return actions
}

func (d *Driver) Reveal(content string) []*Action {
	actions := []*Action{}

	// Get revealable blocks
	var v Vertex
	json.Unmarshal([]byte(content), &v)

	if d.Game.Board[v.Y][v.X].Visited {
		return actions
	}
	revealables := GetRevealables(&v, d.Game.Board)

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
	counter, turn := d.Game.GetCounter(), d.Game.GetTurn()

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

func (d *Driver) startGame() *Action {

	// Init game stat
	d.Game.InitCounter()

	gameStatMsg, _ := json.Marshal(d.GetGameStat())

	return &Action{
		Name:    "gameStat",
		Content: string(gameStatMsg),
	}
}

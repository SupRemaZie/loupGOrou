package application

import (
	"errors"

	"github.com/SupRemaZie/loupGOrou/internal/player"
)

type Phase string

const (
	PhaseDay   Phase = "day"
	PhaseNight Phase = "night"
)

type Game struct {
	ID      string
	Players []*player.Player
	Phase   Phase
	Votes   map[string]int
}

func NewGame(id string, players []*player.Player) *Game {
	return &Game{
		ID:      id,
		Players: players,
		Phase:   PhaseDay,
		Votes:   make(map[string]int),
	}
}

func (g *Game) AddPlayer(p *player.Player) {
	g.Players = append(g.Players, p)
}

func (g *Game) RemovePlayer(playerID string) {
	for i, p := range g.Players {
		if p.ID == playerID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			break
		}
	}
}

func (g *Game) GetPlayer(playerID string) *player.Player {
	for _, p := range g.Players {
		if p.ID == playerID {
			return p
		}
	}
	return nil
}

func (g *Game) AlivePlayers() []*player.Player {
	var alive []*player.Player
	for _, p := range g.Players {
		if p.IsAlive {
			alive = append(alive, p)
		}
	}
	return alive
}

func (g *Game) DeadPlayers() []*player.Player {
	var dead []*player.Player
	for _, p := range g.Players {
		if !p.IsAlive {
			dead = append(dead, p)
		}
	}
	return dead
}

func (g *Game) ResetVotes() {
	g.Votes = make(map[string]int)
	for _, p := range g.Players {
		p.ResetVote()
	}
}

func (g *Game) StartNight() {
	g.Phase = PhaseNight
	g.ResetVotes()
}

func (g *Game) StartDay() {
	g.Phase = PhaseDay
	g.ResetVotes()
}

func (g *Game) Vote(voterID, targetID string) error {
	if g.Phase != PhaseDay {
		return errors.New("le vote est possible uniquement pendant le jour")
	}

	voter := g.GetPlayer(voterID)
	if voter == nil || !voter.IsAlive {
		return errors.New("votant introuvable ou mort")
	}
	if !voter.CanVote() {
		return errors.New("le joueur ne peut pas voter")
	}

	target := g.GetPlayer(targetID)
	if target == nil || !target.IsAlive {
		return errors.New("cible introuvable ou morte")
	}

	if err := voter.Vote(); err != nil {
		return err
	}
	g.Votes[targetID]++
	return nil
}

func (g *Game) ResolveDay() (*player.Player, error) {
	if g.Phase != PhaseDay {
		return nil, errors.New("résolution du jour possible uniquement pendant le jour")
	}

	if len(g.Votes) == 0 {
		return nil, nil
	}

	var loserID string
	maxVotes := 0
	tie := false
	for id, count := range g.Votes {
		if count > maxVotes {
			maxVotes = count
			loserID = id
			tie = false
		} else if count == maxVotes {
			tie = true
		}
	}

	if loserID == "" || tie {
		return nil, nil
	}

	loser := g.GetPlayer(loserID)
	if loser != nil {
		loser.Die()
	}
	return loser, nil
}

func (g *Game) PerformNightAction(actorID, targetID string) error {
	if g.Phase != PhaseNight {
		return errors.New("action de nuit possible uniquement pendant la nuit")
	}

	actor := g.GetPlayer(actorID)
	if actor == nil || !actor.IsAlive {
		return errors.New("acteur introuvable ou mort")
	}

	target := g.GetPlayer(targetID)
	if target == nil || !target.IsAlive {
		return errors.New("cible introuvable ou morte")
	}

	return actor.Role.NightAction(target)
}

func (g *Game) CountFaction(faction string) int {
	count := 0
	for _, p := range g.AlivePlayers() {
		if p.Role.Faction() == faction {
			count++
		}
	}
	return count
}

func (g *Game) IsOver() bool {
	wolves := g.CountFaction("Loup")
	villagers := len(g.AlivePlayers()) - wolves
	return wolves == 0 || wolves >= villagers
}

func (g *Game) Winner() string {
	if !g.IsOver() {
		return "En cours"
	}
	if g.CountFaction("Loup") == 0 {
		return "Villageois"
	}
	return "Loups"
}

package player

import (
	"errors"

	"github.com/SupRemaZie/loupGOrou/internal/role"
)

var (
	ErrPlayerDead   = errors.New("joueur mort, ne peut pas voter")
	ErrAlreadyVoted = errors.New("joueur a déjà voté")
)

type Player struct {
	ID       string
	Name     string
	Role     role.Role
	Mood     string
	IsAlive  bool
	HasVoted bool
}

func NewPlayer(id, name string, role role.Role) *Player {
	return &Player{
		ID:       id,
		Name:     name,
		Role:     role,
		Mood:     "neutral",
		IsAlive:  true,
		HasVoted: false,
	}
}

func (p *Player) Die() {
	p.IsAlive = false
}

func (p *Player) ResetVote() {
	p.HasVoted = false
}

func (p *Player) CanVote() bool {
	return p.IsAlive && !p.HasVoted
}

func (p *Player) Vote() error {
	if !p.IsAlive {
		return ErrPlayerDead
	}
	if p.HasVoted {
		return ErrAlreadyVoted
	}
	p.HasVoted = true
	return nil
}

func (p *Player) GetFaction() string {
	return p.Role.Faction()
}

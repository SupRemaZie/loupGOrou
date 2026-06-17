package player

import (
	"github.com/SupRemaZie/loupGOrou/internal/role"
)

type Player struct {
	Name    string
	Role    role.Role
	IsAlive bool
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:    name,
		IsAlive: true,
	}
}

func (p *Player) Die() {
	p.IsAlive = false
}

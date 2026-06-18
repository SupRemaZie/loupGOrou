package phase

import (
	"fmt"

	"github.com/SupRemaZie/loupGOrou/internal/action"
	"github.com/SupRemaZie/loupGOrou/internal/player"
)

type nightPhase struct {
	alivePlayers []*player.Player
}

func NewNightPhase(alivePlayers []*player.Player) *nightPhase {
	return &nightPhase{
		alivePlayers: alivePlayers,
	}
}

func (np *nightPhase) Start() []*player.Player {
	fmt.Println("\n🌙 La nuit tombe sur le village...")

	ctx := &nightContext{
		alive:   np.alivePlayers,
		pending: make(map[*player.Player]bool),
	}

	for _, a := range action.BuildNightActions(np.alivePlayers) {
		a.Resolve(ctx)
	}

	for _, p := range np.alivePlayers {
		p.Role.ResetNight()
	}

	victims := ctx.Victims()
	for _, v := range victims {
		v.Die()
	}
	return victims
}

type nightContext struct {
	alive   []*player.Player
	pending map[*player.Player]bool
}

func (c *nightContext) Kill(p *player.Player) {
	if p != nil {
		c.pending[p] = true
	}
}

func (c *nightContext) Save(p *player.Player) {
	delete(c.pending, p)
}

func (c *nightContext) Victims() []*player.Player {
	var res []*player.Player
	for _, p := range c.alive {
		if c.pending[p] {
			res = append(res, p)
		}
	}
	return res
}

func (c *nightContext) GetAlive() []*player.Player {
	return c.alive
}

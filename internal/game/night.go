package game

import (
	"fmt"
	"math/rand"

	"github.com/SupRemaZie/loupGOrou/internal/player"
)

type nightPhase struct {
	alivePlayers []*player.Player
}

func newNightPhase(alivePlayers []*player.Player) *nightPhase {
	return &nightPhase{
		alivePlayers: alivePlayers,
	}
}

func (np *nightPhase) start() []*player.Player {
	fmt.Println("\n🌙 La nuit tombe sur le village...")

	ctx := &nightContext{
		alive:   np.alivePlayers,
		pending: make(map[*player.Player]bool),
	}

	for _, action := range buildNightActions(np.alivePlayers) {
		action.Resolve(ctx)
	}

	for _, p := range np.alivePlayers {
		p.Role.ResetNight()
	}

	victims := ctx.victims()
	for _, v := range victims {
		v.Die()
	}
	return victims
}

type nightContext struct {
	alive   []*player.Player
	pending map[*player.Player]bool
}

func (c *nightContext) kill(p *player.Player) {
	if p != nil {
		c.pending[p] = true
	}
}

func (c *nightContext) save(p *player.Player) {
	delete(c.pending, p)
}

func (c *nightContext) victims() []*player.Player {
	var res []*player.Player
	for _, p := range c.alive {
		if c.pending[p] {
			res = append(res, p)
		}
	}
	return res
}

func pickVictim(votes map[*player.Player]int) *player.Player {
	maxVotes := 0
	for _, count := range votes {
		if count > maxVotes {
			maxVotes = count
		}
	}
	if maxVotes == 0 {
		return nil
	}

	var top []*player.Player
	for p, count := range votes {
		if count == maxVotes {
			top = append(top, p)
		}
	}
	return top[rand.Intn(len(top))]
}

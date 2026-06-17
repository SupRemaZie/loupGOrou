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

	wolves := faction(np.alivePlayers, "Loup")
	if len(wolves) == 0 {
		return nil
	}

	prey := otherFactions(wolves[0], np.alivePlayers)

	votes := make(map[*player.Player]int)
	for _, wolf := range wolves {
		if !wolf.IsAlive || !wolf.Role.CanAct() {
			continue
		}
		fmt.Printf("\n%s se réveille (%s).\n", wolf.Name, wolf.Role.Name())
		if target := promptTarget(wolf, prey); target != nil {
			votes[target]++
		}
	}

	for _, p := range np.alivePlayers {
		p.Role.ResetNight()
	}

	victim := pickVictim(votes)
	if victim == nil {
		return nil
	}
	victim.Die()
	return []*player.Player{victim}
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

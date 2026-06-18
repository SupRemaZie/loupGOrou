package game

import (
	"fmt"

	"github.com/SupRemaZie/loupGOrou/internal/phase"
	"github.com/SupRemaZie/loupGOrou/internal/player"
)

type round struct {
	alivePlayers []*player.Player
}

func newRound(alivePlayers []*player.Player) *round {
	return &round{
		alivePlayers: alivePlayers,
	}
}

func (r *round) start() {
	night := phase.NewNightPhase(r.alivePlayers)
	for _, victim := range night.Start() {
		fmt.Printf("Cette nuit, %s a été dévoré...\n", victim.Name)
	}

	day := phase.NewDayPhase(alive(r.alivePlayers))
	if suspect := day.Start(); suspect != nil {
		fmt.Printf("Le village a éliminé %s (%s).\n", suspect.Name, suspect.Role.Name())
	}
}

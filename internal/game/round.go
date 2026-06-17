package game

import (
	"fmt"

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
	night := newNightPhase(r.alivePlayers)
	for _, victim := range night.start() {
		fmt.Printf("Cette nuit, %s a été dévoré...\n", victim.Name)
	}

	day := newDayPhase(alive(r.alivePlayers))
	if suspect := day.start(); suspect != nil {
		fmt.Printf("Le village a éliminé %s (%s).\n", suspect.Name, suspect.Role.Name())
	}
}

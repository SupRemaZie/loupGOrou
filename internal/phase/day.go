package phase

import (
	"fmt"

	"github.com/SupRemaZie/loupGOrou/internal/action"
	"github.com/SupRemaZie/loupGOrou/internal/player"
)

type dayPhase struct {
	alivePlayers []*player.Player
}

func NewDayPhase(alivePlayers []*player.Player) *dayPhase {
	return &dayPhase{
		alivePlayers: alivePlayers,
	}
}

func (dp *dayPhase) Start() *player.Player {
	fmt.Println("\n☀️ Le jour se lève. Le village doit voter pour un suspect.")

	votes := make(map[*player.Player]int)
	for _, voter := range dp.alivePlayers {
		if !voter.IsAlive {
			continue
		}
		fmt.Printf("\n%s, à qui voulez-vous donner votre voix ?\n", voter.Name)
		if target := action.PromptTarget(voter, dp.alivePlayers); target != nil {
			votes[target]++
		}
	}

	var suspect *player.Player
	maxVotes := 0
	tie := false
	for p, count := range votes {
		switch {
		case count > maxVotes:
			maxVotes, suspect, tie = count, p, false
		case count == maxVotes:
			tie = true
		}
	}

	if suspect == nil || tie {
		fmt.Println("Aucun consensus : personne n'est éliminé.")
		return nil
	}

	suspect.Die()
	return suspect
}

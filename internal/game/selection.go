package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SupRemaZie/loupGOrou/internal/console"
	"github.com/SupRemaZie/loupGOrou/internal/player"
)

func alive(players []*player.Player) []*player.Player {
	var res []*player.Player
	for _, p := range players {
		if p.IsAlive {
			res = append(res, p)
		}
	}
	return res
}

func faction(players []*player.Player, name string) []*player.Player {
	var res []*player.Player
	for _, p := range players {
		if p.IsAlive && p.Role.Faction() == name {
			res = append(res, p)
		}
	}
	return res
}

func otherFactions(actor *player.Player, players []*player.Player) []*player.Player {
	var res []*player.Player
	for _, p := range players {
		if p.IsAlive && p.Role.Faction() != actor.Role.Faction() {
			res = append(res, p)
		}
	}
	return res
}

func promptTarget(actor *player.Player, candidates []*player.Player) *player.Player {
	var choices []*player.Player
	for _, p := range candidates {
		if p.IsAlive && p != actor {
			choices = append(choices, p)
		}
	}
	if len(choices) == 0 {
		return nil
	}

	for i, p := range choices {
		fmt.Printf("  %d. %s\n", i+1, p.Name)
	}

	for {
		fmt.Printf("%s, choisissez une cible (numéro ou nom) : ", actor.Name)
		input := console.ReadLine()
		if input == "" {
			continue
		}
		if n, err := strconv.Atoi(input); err == nil {
			if n >= 1 && n <= len(choices) {
				return choices[n-1]
			}
			fmt.Println("Numéro hors limites, réessayez.")
			continue
		}
		for _, p := range choices {
			if strings.EqualFold(p.Name, input) {
				return p
			}
		}
		fmt.Println("Choix invalide, réessayez.")
	}
}

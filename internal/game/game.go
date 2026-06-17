package game

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/SupRemaZie/loupGOrou/internal/console"
	"github.com/SupRemaZie/loupGOrou/internal/player"
	"github.com/SupRemaZie/loupGOrou/internal/role"
)

type Game struct {
	Players    []*player.Player
	NbWerewolf int
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Start(nbPlayers int, nbWerewolves int) error {
	if nbPlayers <= 0 {
		return errors.New("impossible de démarrer sans joueurs")
	}
	for nbWerewolves <= 0 || nbWerewolves*2 >= nbPlayers {
		fmt.Print("Le nombre de loups est trop élevé. Réessayez : ")
		nbWerewolves = console.ReadInt()
	}
	g.NbWerewolf = nbWerewolves

	for i := range nbPlayers {
		fmt.Printf("Entrez le nom du joueur %d : ", i+1)
		g.AddPlayer(player.NewPlayer(console.ReadLine()))
	}

	rand.Shuffle(len(g.Players), func(i, j int) {
		g.Players[i], g.Players[j] = g.Players[j], g.Players[i]
	})

	for i := 0; i < g.NbWerewolf; i++ {
		g.Players[i].Role = role.NewWerewolf()
	}

	// Rôles spéciaux du village, attribués si l'effectif le permet.
	next := g.NbWerewolf
	if next < len(g.Players) {
		g.Players[next].Role = role.NewSeer()
		next++
	}
	if next < len(g.Players) {
		g.Players[next].Role = role.NewWitch()
		next++
	}

	for i := next; i < len(g.Players); i++ {
		g.Players[i].Role = role.NewVillager()
	}

	roundNumber := 1
	for !g.IsOver() {
		fmt.Printf("\n===== Round %d =====\n", roundNumber)
		newRound(g.AlivePlayers()).start()
		roundNumber++
	}

	fmt.Printf("\nFin de la partie ! Vainqueur : %s\n", g.Winner())
	return nil
}

func (g *Game) AddPlayer(p *player.Player) {
	g.Players = append(g.Players, p)
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

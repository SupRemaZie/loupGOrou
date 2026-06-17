package main

import (
	"fmt"

	"github.com/SupRemaZie/loupGOrou/internal/console"
	"github.com/SupRemaZie/loupGOrou/internal/game"
)

func main() {
	fmt.Println("Bienvenue dans le jeu du Loup Garou !")

	fmt.Println("Entrez le nombre de joueurs :")
	nbPlayers := console.ReadInt()
	fmt.Println("Entrez le nombre de loups :")
	nbWerewolves := console.ReadInt()

	g := game.NewGame()
	if err := g.Start(nbPlayers, nbWerewolves); err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	for i, p := range g.Players {
		fmt.Printf("Joueur %d : %s, Rôle : %s, Faction : %s\n", i+1, p.Name, p.Role.String(), p.Role.Faction())
	}
}

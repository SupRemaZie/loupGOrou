package main

import (
	"fmt"

	"github.com/SupRemaZie/loupGOrou/internal/application"
	"github.com/SupRemaZie/loupGOrou/internal/player"
	"github.com/SupRemaZie/loupGOrou/internal/role"
)

func main() {
	player1 := player.NewPlayer("1", "Alice", role.NewWerewolf())
	player2 := player.NewPlayer("2", "Bob", role.NewVillager())
	player3 := player.NewPlayer("3", "Charlie", role.NewVillager())

	players := []*player.Player{player1, player2, player3}
	game := application.NewGame("game-1", players)

	fmt.Printf("Game created: %s, phase=%s, players=%d\n", game.ID, game.Phase, len(game.Players))
}

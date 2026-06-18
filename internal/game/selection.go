package game

import "github.com/SupRemaZie/loupGOrou/internal/player"

func alive(players []*player.Player) []*player.Player {
	var res []*player.Player
	for _, p := range players {
		if p.IsAlive {
			res = append(res, p)
		}
	}
	return res
}

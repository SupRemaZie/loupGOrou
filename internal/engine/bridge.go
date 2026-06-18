package engine

import (
	"github.com/SupRemaZie/loupGOrou/internal/player"
	"github.com/SupRemaZie/loupGOrou/internal/role"
)

func hydrate(state GameState) []*player.Player {
	players := make([]*player.Player, 0, len(state.Players))
	for _, ps := range state.Players {
		p := player.NewPlayer(ps.Name)
		p.IsAlive = ps.IsAlive
		p.Role = role.FromName(ps.Role)

		if w, ok := p.Role.(*role.Witch); ok {
			if !ps.HasHeal {
				w.UseHeal()
			}
			if !ps.HasPoison {
				w.UsePoison()
			}
		}
		players = append(players, p)
	}
	return players
}

func dehydrate(gameID string, players []*player.Player, round int, phase Phase) GameState {
	states := make([]PlayerState, len(players))
	for i, p := range players {
		ps := PlayerState{
			ID:      p.Name,
			Name:    p.Name,
			IsAlive: p.IsAlive,
			Role:    p.Role.Name(),
			Faction: p.Role.Faction(),
		}
		if w, ok := p.Role.(*role.Witch); ok {
			ps.HasHeal = w.HasHeal()
			ps.HasPoison = w.HasPoison()
		}
		states[i] = ps
	}
	return GameState{ID: gameID, Round: round, Phase: phase, Players: states}
}

func findDecision(decisions []Decision, kind DecisionKind, actorID string) *Decision {
	for i := range decisions {
		if decisions[i].Kind == kind && decisions[i].ActorID == actorID {
			return &decisions[i]
		}
	}
	return nil
}

func findPlayer(players []*player.Player, id string) *player.Player {
	for _, p := range players {
		if p.Name == id && p.IsAlive {
			return p
		}
	}
	return nil
}

func alivePlayers(players []*player.Player) []*player.Player {
	var alive []*player.Player
	for _, p := range players {
		if p.IsAlive {
			alive = append(alive, p)
		}
	}
	return alive
}

func toPublicStates(players []*player.Player) []PlayerState {
	states := make([]PlayerState, 0, len(players))
	for _, p := range players {
		if p.IsAlive {
			states = append(states, PlayerState{
				ID:      p.Name,
				Name:    p.Name,
				IsAlive: true,
			})
		}
	}
	return states
}

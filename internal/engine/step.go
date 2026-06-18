package engine

import (
	"math/rand"

	"github.com/SupRemaZie/loupGOrou/internal/player"
	"github.com/SupRemaZie/loupGOrou/internal/role"
)

// ComputeRequired retourne toutes les décisions (obligatoires et optionnelles) pour la phase courante.
func ComputeRequired(state GameState) []RequiredDecision {
	players := hydrate(state)
	return computeRequired(players, state)
}

// Step fait avancer le jeu d'une demi-phase.
// Si des décisions obligatoires manquent, retourne Pending sans rien modifier.
func Step(state GameState, decisions []Decision) StepResult {
	players := hydrate(state)
	required := computeRequired(players, state)
	if missing := filterMissing(required, decisions); len(missing) > 0 {
		return StepResult{State: state, Pending: missing}
	}

	var events []Event
	switch state.Phase {
	case PhaseNight:
		events, state = executeNight(state, players, decisions)
	case PhaseNightWitch:
		events, state = executeNightWitch(state, players, decisions)
	case PhaseDay:
		events, state = executeDay(state, players, decisions)
	}
	return StepResult{State: state, Events: events}
}

// computeRequired liste les décisions nécessaires pour la phase courante.
func computeRequired(players []*player.Player, state GameState) []RequiredDecision {
	var req []RequiredDecision
	alive := alivePlayers(players)

	switch state.Phase {
	case PhaseNight:
		wolves := filterFaction(alive, "Loup")
		prey := filterNotFaction(alive, "Loup")
		for _, wolf := range wolves {
			req = append(req, RequiredDecision{
				Kind:       DecisionWerewolfAttack,
				ActorID:    wolf.Name,
				Candidates: toPublicStates(prey),
			})
		}
		if seer := findByRole(alive, "Voyante"); seer != nil {
			req = append(req, RequiredDecision{
				Kind:       DecisionSeerInvestigate,
				ActorID:    seer.Name,
				Candidates: toPublicStates(filterOut(alive, seer)),
			})
		}

	case PhaseNightWitch:
		witch := findByRole(alive, "Sorcière")
		if witch == nil {
			break
		}
		w := witch.Role.(*role.Witch)
		if w.HasHeal() && state.Victim != "" {
			if victim := findPlayer(players, state.Victim); victim != nil {
				req = append(req, RequiredDecision{
					Kind:       DecisionWitchSave,
					ActorID:    witch.Name,
					Candidates: []PlayerState{{ID: victim.Name, Name: victim.Name, IsAlive: true}},
					Optional:   true,
				})
			}
		}
		if w.HasPoison() {
			req = append(req, RequiredDecision{
				Kind:       DecisionWitchPoison,
				ActorID:    witch.Name,
				Candidates: toPublicStates(filterOut(alive, witch)),
				Optional:   true,
			})
		}

	case PhaseDay:
		for _, p := range alive {
			req = append(req, RequiredDecision{
				Kind:       DecisionVote,
				ActorID:    p.Name,
				Candidates: toPublicStates(filterOut(alive, p)),
			})
		}
	}
	return req
}

// filterMissing retourne les décisions obligatoires non fournies.
func filterMissing(required []RequiredDecision, decisions []Decision) []RequiredDecision {
	var missing []RequiredDecision
	for _, req := range required {
		if req.Optional {
			continue
		}
		if findDecision(decisions, req.Kind, req.ActorID) == nil {
			missing = append(missing, req)
		}
	}
	return missing
}

func executeNight(state GameState, players []*player.Player, decisions []Decision) ([]Event, GameState) {
	alive := alivePlayers(players)
	var events []Event

	// Vote des loups
	wolves := filterFaction(alive, "Loup")
	prey := filterNotFaction(alive, "Loup")
	victim := wolfVictim(wolves, prey, decisions)
	victimID := ""
	if victim != nil {
		victimID = victim.Name
	}

	// Enquête de la voyante
	if seer := findByRole(alive, "Voyante"); seer != nil {
		if d := findDecision(decisions, DecisionSeerInvestigate, seer.Name); d != nil && d.TargetID != "" {
			if target := findPlayer(alive, d.TargetID); target != nil {
				events = append(events, Event{
					Kind:     EventRevealed,
					PlayerID: target.Name,
					Detail:   target.Role.Name(),
				})
			}
		}
	}

	// Si la sorcière peut encore agir : suspend les kills et passe en PhaseNightWitch
	if witch := findByRole(alive, "Sorcière"); witch != nil {
		if w, ok := witch.Role.(*role.Witch); ok && w.CanAct() {
			next := dehydrate(state.ID, players, state.Round, PhaseNightWitch)
			next.Victim = victimID
			return events, next
		}
	}

	// Pas de sorcière : applique le kill immédiatement
	if victim != nil {
		victim.Die()
		events = append(events, Event{Kind: EventKilled, PlayerID: victim.Name})
	}
	next := dehydrate(state.ID, players, state.Round, PhaseDay)
	next.Result = checkResult(alivePlayers(players))
	return events, next
}

func executeNightWitch(state GameState, players []*player.Player, decisions []Decision) ([]Event, GameState) {
	alive := alivePlayers(players)
	var events []Event

	var witch *role.Witch
	var witchPlayer *player.Player
	if wp := findByRole(alive, "Sorcière"); wp != nil {
		if w, ok := wp.Role.(*role.Witch); ok {
			witch = w
			witchPlayer = wp
		}
	}

	// Potion de soin
	victimSaved := false
	if witch != nil && witch.HasHeal() && state.Victim != "" {
		if d := findDecision(decisions, DecisionWitchSave, witchPlayer.Name); d != nil && d.TargetID != "" {
			witch.UseHeal()
			victimSaved = true
			events = append(events, Event{Kind: EventSaved, PlayerID: state.Victim})
		}
	}

	// Potion de poison
	if witch != nil && witch.HasPoison() {
		if d := findDecision(decisions, DecisionWitchPoison, witchPlayer.Name); d != nil && d.TargetID != "" {
			if target := findPlayer(alive, d.TargetID); target != nil {
				witch.UsePoison()
				target.Die()
				events = append(events, Event{Kind: EventKilled, PlayerID: target.Name})
			}
		}
	}

	// Applique le kill des loups si la victime n'a pas été sauvée
	if !victimSaved && state.Victim != "" {
		if victim := findPlayer(players, state.Victim); victim != nil {
			victim.Die()
			events = append(events, Event{Kind: EventKilled, PlayerID: victim.Name})
		}
	}

	next := dehydrate(state.ID, players, state.Round, PhaseDay)
	next.Result = checkResult(alivePlayers(players))
	return events, next
}

func executeDay(state GameState, players []*player.Player, decisions []Decision) ([]Event, GameState) {
	alive := alivePlayers(players)
	var events []Event

	votes := make(map[*player.Player]int)
	for _, p := range alive {
		if d := findDecision(decisions, DecisionVote, p.Name); d != nil && d.TargetID != "" {
			if target := findPlayer(alive, d.TargetID); target != nil {
				votes[target]++
			}
		}
	}

	if suspect := pickByVotes(votes); suspect != nil {
		suspect.Die()
		events = append(events, Event{
			Kind:     EventEliminated,
			PlayerID: suspect.Name,
			Detail:   suspect.Role.Name(),
		})
	} else {
		events = append(events, Event{Kind: EventNoConsensus})
	}

	next := dehydrate(state.ID, players, state.Round+1, PhaseNight)
	next.Result = checkResult(alivePlayers(players))
	return events, next
}

// --- helpers ---

func wolfVictim(wolves, prey []*player.Player, decisions []Decision) *player.Player {
	votes := make(map[*player.Player]int)
	for _, wolf := range wolves {
		if !wolf.Role.CanAct() {
			continue
		}
		d := findDecision(decisions, DecisionWerewolfAttack, wolf.Name)
		if d == nil || d.TargetID == "" {
			continue
		}
		if t := findPlayer(prey, d.TargetID); t != nil {
			votes[t]++
		}
	}
	return pickByVotes(votes)
}

// pickByVotes retourne le joueur avec le plus de votes, nil en cas d'égalité.
func pickByVotes(votes map[*player.Player]int) *player.Player {
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
	if len(top) > 1 {
		return nil // égalité = personne n'est éliminé
	}
	return top[rand.Intn(len(top))]
}

func filterFaction(players []*player.Player, faction string) []*player.Player {
	var res []*player.Player
	for _, p := range players {
		if p.Role.Faction() == faction {
			res = append(res, p)
		}
	}
	return res
}

func filterNotFaction(players []*player.Player, faction string) []*player.Player {
	var res []*player.Player
	for _, p := range players {
		if p.Role.Faction() != faction {
			res = append(res, p)
		}
	}
	return res
}

func filterOut(players []*player.Player, exclude *player.Player) []*player.Player {
	var res []*player.Player
	for _, p := range players {
		if p != exclude {
			res = append(res, p)
		}
	}
	return res
}

func findByRole(players []*player.Player, roleName string) *player.Player {
	for _, p := range players {
		if p.Role.Name() == roleName {
			return p
		}
	}
	return nil
}

func checkResult(alive []*player.Player) string {
	wolves := 0
	for _, p := range alive {
		if p.Role.Faction() == "Loup" {
			wolves++
		}
	}
	if wolves == 0 {
		return "Villageois"
	}
	if wolves >= len(alive)-wolves {
		return "Loups"
	}
	return ""
}

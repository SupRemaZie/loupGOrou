package action

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/SupRemaZie/loupGOrou/internal/console"
	"github.com/SupRemaZie/loupGOrou/internal/player"
	"github.com/SupRemaZie/loupGOrou/internal/role"
)

type NightContext interface {
	Kill(p *player.Player)
	Save(p *player.Player)
	Victims() []*player.Player
	GetAlive() []*player.Player
}

type NightAction interface {
	Priority() int
	Resolve(ctx NightContext)
}

func BuildNightActions(alive []*player.Player) []NightAction {
	var actions []NightAction

	if wolves := faction(alive, "Loup"); len(wolves) > 0 {
		actions = append(actions, &wolfPackAction{wolves: wolves})
	}

	for _, p := range alive {
		switch p.Role.Name() {
		case "Voyante":
			actions = append(actions, &seerAction{self: p})
		case "Sorcière":
			actions = append(actions, &witchAction{self: p})
		}
	}

	sort.SliceStable(actions, func(i, j int) bool {
		return actions[i].Priority() < actions[j].Priority()
	})
	return actions
}

func pickVictim(votes map[*player.Player]int) *player.Player {
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
	return top[rand.Intn(len(top))]
}

type seerAction struct{ self *player.Player }

func (a *seerAction) Priority() int { return 10 }

func (a *seerAction) Resolve(ctx NightContext) {
	if !a.self.IsAlive {
		return
	}
	fmt.Printf("\n%s se réveille (Voyante).\n", a.self.Name)
	if target := PromptTarget(a.self, ctx.GetAlive()); target != nil {
		fmt.Printf("🔮 %s est : %s\n", target.Name, target.Role.Name())
	}
}

type wolfPackAction struct{ wolves []*player.Player }

func (a *wolfPackAction) Priority() int { return 20 }

func (a *wolfPackAction) Resolve(ctx NightContext) {
	prey := otherFactions(a.wolves[0], ctx.GetAlive())

	votes := make(map[*player.Player]int)
	for _, wolf := range a.wolves {
		if !wolf.IsAlive || !wolf.Role.CanAct() {
			continue
		}
		fmt.Printf("\n%s se réveille (%s).\n", wolf.Name, wolf.Role.Name())
		if target := PromptTarget(wolf, prey); target != nil {
			votes[target]++
		}
	}
	ctx.Kill(pickVictim(votes))
}

type witchAction struct{ self *player.Player }

func (a *witchAction) Priority() int { return 30 }

func (a *witchAction) Resolve(ctx NightContext) {
	if !a.self.IsAlive {
		return
	}
	witch, ok := a.self.Role.(*role.Witch)
	if !ok {
		return
	}
	fmt.Printf("\n%s se réveille (Sorcière).\n", a.self.Name)

	if witch.HasHeal() {
		if victims := ctx.Victims(); len(victims) > 0 {
			v := victims[0]
			fmt.Printf("Cette nuit, %s va mourir. Le sauver ? (o/n) : ", v.Name)
			if askYesNo() {
				ctx.Save(v)
				witch.UseHeal()
			}
		}
	}

	if witch.HasPoison() {
		fmt.Print("Empoisonner quelqu'un ? (o/n) : ")
		if askYesNo() {
			if target := PromptTarget(a.self, ctx.GetAlive()); target != nil {
				ctx.Kill(target)
				witch.UsePoison()
			}
		}
	}
}

func askYesNo() bool {
	return strings.EqualFold(console.ReadLine(), "o")
}

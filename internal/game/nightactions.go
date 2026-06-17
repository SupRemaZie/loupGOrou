package game

import (
	"fmt"
	"sort"
	"strings"

	"github.com/SupRemaZie/loupGOrou/internal/console"
	"github.com/SupRemaZie/loupGOrou/internal/player"
	"github.com/SupRemaZie/loupGOrou/internal/role"
)

type NightAction interface {
	Priority() int
	Resolve(ctx *nightContext)
}

func buildNightActions(alive []*player.Player) []NightAction {
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

type seerAction struct{ self *player.Player }

func (a *seerAction) Priority() int { return 10 }

func (a *seerAction) Resolve(ctx *nightContext) {
	if !a.self.IsAlive {
		return
	}
	fmt.Printf("\n%s se réveille (Voyante).\n", a.self.Name)
	if target := promptTarget(a.self, ctx.alive); target != nil {
		fmt.Printf("🔮 %s est : %s\n", target.Name, target.Role.Name())
	}
}

type wolfPackAction struct{ wolves []*player.Player }

func (a *wolfPackAction) Priority() int { return 20 }

func (a *wolfPackAction) Resolve(ctx *nightContext) {
	prey := otherFactions(a.wolves[0], ctx.alive)

	votes := make(map[*player.Player]int)
	for _, wolf := range a.wolves {
		if !wolf.IsAlive || !wolf.Role.CanAct() {
			continue
		}
		fmt.Printf("\n%s se réveille (%s).\n", wolf.Name, wolf.Role.Name())
		if target := promptTarget(wolf, prey); target != nil {
			votes[target]++
		}
	}
	ctx.kill(pickVictim(votes))
}

type witchAction struct{ self *player.Player }

func (a *witchAction) Priority() int { return 30 }

func (a *witchAction) Resolve(ctx *nightContext) {
	if !a.self.IsAlive {
		return
	}
	witch, ok := a.self.Role.(*role.Witch)
	if !ok {
		return
	}
	fmt.Printf("\n%s se réveille (Sorcière).\n", a.self.Name)

	if witch.HasHeal() {
		if victims := ctx.victims(); len(victims) > 0 {
			v := victims[0]
			fmt.Printf("Cette nuit, %s va mourir. Le sauver ? (o/n) : ", v.Name)
			if askYesNo() {
				ctx.save(v)
				witch.UseHeal()
			}
		}
	}

	if witch.HasPoison() {
		fmt.Print("Empoisonner quelqu'un ? (o/n) : ")
		if askYesNo() {
			if target := promptTarget(a.self, ctx.alive); target != nil {
				ctx.kill(target)
				witch.UsePoison()
			}
		}
	}
}

func askYesNo() bool {
	return strings.EqualFold(console.ReadLine(), "o")
}

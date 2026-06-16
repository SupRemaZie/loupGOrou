package role

type Villager struct {
}

func (v Villager) Name() string        { return "Villageois" }
func (v Villager) Faction() string     { return "civil" }
func (v Villager) NightActions() error { return nil }
func (v Villager) CanAct() bool        { return false }

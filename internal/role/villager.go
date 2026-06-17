package role

type Villager struct {
}

func NewVillager() Villager {
	return Villager{}
}

func (v Villager) Name() string    { return "Villageois" }
func (v Villager) Faction() string { return "civil" }
func (v Villager) CanAct() bool    { return false }
func (v Villager) NightAction(target AttackTarget) error {
	return nil
}
func (v Villager) ResetNight() {}

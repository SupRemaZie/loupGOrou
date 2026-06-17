package role

type AttackTarget interface {
	Die()
}

type Role interface {
	Name() string
	Faction() string
	CanAct() bool
	NightAction(target AttackTarget) error
	ResetNight()
}

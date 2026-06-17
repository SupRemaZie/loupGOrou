package role

type Role interface {
	Name() string
	Faction() string
	CanAct() bool
	NightActions() error
}

package role

type Werewolf struct{}

func NewWerewolf() Werewolf { return Werewolf{} }

func (w Werewolf) Name() string    { return "Loup Garou" }
func (w Werewolf) Faction() string { return "Loup" }
func (w Werewolf) CanAct() bool    { return true }
func (w Werewolf) ResetNight()     {}
func (w Werewolf) String() string  { return "🐺 " + w.Name() }

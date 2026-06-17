package role

type Seer struct{}

func NewSeer() Seer { return Seer{} }

func (s Seer) Name() string    { return "Voyante" }
func (s Seer) Faction() string { return "civil" }
func (s Seer) CanAct() bool    { return true }
func (s Seer) ResetNight()     {}
func (s Seer) String() string  { return "🔮 " + s.Name() }

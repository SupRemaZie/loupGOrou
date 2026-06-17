package role

type Witch struct {
	healUsed   bool
	poisonUsed bool
}

func NewWitch() *Witch { return &Witch{} }

func (w *Witch) Name() string    { return "Sorcière" }
func (w *Witch) Faction() string { return "civil" }
func (w *Witch) CanAct() bool    { return !w.healUsed || !w.poisonUsed }
func (w *Witch) ResetNight()     {}
func (w *Witch) String() string  { return "🧪 " + w.Name() }

func (w *Witch) HasHeal() bool   { return !w.healUsed }
func (w *Witch) HasPoison() bool { return !w.poisonUsed }
func (w *Witch) UseHeal()        { w.healUsed = true }
func (w *Witch) UsePoison()      { w.poisonUsed = true }

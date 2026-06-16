package role

type Werewolf struct {
	hasActedThisNight bool
}

func (w Werewolf) Name() string    { return "Loup Garou" }
func (w Werewolf) Faction() string { return "Loup" }
func (w Werewolf) CanAct() bool    { return !w.hasActedThisNight }

func (w *Werewolf) NightActions() error {
	w.hasActedThisNight = true
	return nil
}

func (w *Werewolf) ResetNight() {
	w.hasActedThisNight = false
}

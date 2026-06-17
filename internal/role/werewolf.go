package role

import "errors"

type Werewolf struct {
	hasActedThisNight bool
}

func NewWerewolf() *Werewolf {
	return &Werewolf{hasActedThisNight: false}
}

func (w Werewolf) Name() string    { return "Loup Garou" }
func (w Werewolf) Faction() string { return "Loup" }
func (w Werewolf) CanAct() bool    { return !w.hasActedThisNight }

func (w *Werewolf) NightAction(target AttackTarget) error {
	if !w.CanAct() {
		return errors.New("le loup a déjà agi cette nuit")
	}
	if target == nil {
		return errors.New("cible invalide")
	}
	target.Die()
	w.hasActedThisNight = true
	return nil
}

func (w *Werewolf) ResetNight() {
	w.hasActedThisNight = false
}

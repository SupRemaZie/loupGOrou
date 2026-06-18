package role

func FromName(name string) Role {
	switch name {
	case "Loup Garou":
		return NewWerewolf()
	case "Voyante":
		return NewSeer()
	case "Sorcière":
		return NewWitch()
	default:
		return NewVillager()
	}
}

package models

type EntitiesByType struct {
	Senders []*Entity
	All     []*Entity
}

func NewEntitiesByType(entities []*Entity) *EntitiesByType {
	entitiesByType := &EntitiesByType{
		Senders: []*Entity{},
		All:     []*Entity{},
	}

	for _, e := range entities {
		if e.UserType == EntityUserTypes.PRODUTOR_RURAL() {
			entitiesByType.Senders = append(entitiesByType.Senders, e)
		}
		entitiesByType.All = append(entitiesByType.All, e)
	}

	return entitiesByType
}

package model

import "github.com/NumberMan1/common/summer/vector3"

type Actor struct {
	*Entity
	Name  string
	Level int
	Speed int
}

func NewActor(entityId int, position, direction vector3.Vector3) *Actor {
	return &Actor{Entity: NewEntity(entityId, position, direction)}
}

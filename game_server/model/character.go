package model

import (
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/vector3"
)

// Character 角色
type Character struct {
	*Actor
	//当前角色的客户端连接
	Conn network.Connection
}

func NewCharacter(id int, position, direction vector3.Vector3) *Character {
	return &Character{Actor: NewActor(id, position, direction)}
}

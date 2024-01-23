package model

import (
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/network"
)

var ( //当前场景中全部的角色
	characterDict = map[int]*Character{}
	connCharacter = map[network.Connection]*Character{}
)

type Space struct {
	Id   int
	Name string
}

// CharacterJoin 角色进入场景
func (s Space) CharacterJoin(conn network.Connection, character *Character) {
	logger.SLCInfo("角色进入场景:%d", character.EntityId())

}

package model

import (
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/network"
	pt "github.com/NumberMan1/common/summer/protocol/gen/proto"
)

type Space struct {
	Id   int
	Name string
	//当前场景中全部的角色
	characterDict map[int]*Character
	connCharacter map[network.Connection]*Character
}

func NewSpace(id int, name string) *Space {
	return &Space{Id: id, Name: name, characterDict: map[int]*Character{}, connCharacter: map[network.Connection]*Character{}}
}

// CharacterJoin 角色进入场景
func (s *Space) CharacterJoin(conn network.Connection, character *Character) {
	logger.SLCInfo("角色进入场景:%d", character.EntityId())
	conn.Set("Character", character) //把角色存入连接
	character.SpaceId = s.Id
	character.Conn = conn
	s.characterDict[character.EntityId()] = character
	_, ok := s.connCharacter[conn]
	if ok == false {
		s.connCharacter[conn] = character
	}
	conn.Set("Space", s) //把场景存入连接
	//把新进入的角色广播给其他玩家
	response := &pt.SpaceCharactersEnterResponse{
		SpaceId:    int32(s.Id),
		EntityList: make([]*pt.NEntity, 0),
	}
	response.EntityList = append(response.EntityList, character.GetData())
	for _, v := range s.characterDict {
		if v.Conn != conn {
			v.Conn.Send(response)
		}
	}
	//新上线的角色需要获取全部角色
	for _, v := range s.characterDict {
		if v.Conn == conn {
			continue
		}
		response.EntityList = make([]*pt.NEntity, 0)
		response.EntityList = append(response.EntityList, v.GetData())
		conn.Send(response)
	}
}

// CharacterLeave 角色离开地图
// 客户端离线、切换地图
func (s *Space) CharacterLeave(conn network.Connection, character *Character) {
	logger.SLCInfo("角色离开场景:%d", character.EntityId())
	conn.Set("Space", nil)
	delete(s.characterDict, character.EntityId()) //取消conn的场景记录
	response := &pt.SpaceCharactersEnterResponse{
		SpaceId:    int32(s.Id),
		EntityList: nil,
	}
	for _, v := range s.characterDict {
		v.Conn.Send(response)
	}
}

// UpdateEntity 广播更新Entity信息
func (s *Space) UpdateEntity(sync *pt.NEntitySync) {
	logger.SLCInfo("UpdateEntity %s", sync.String())
	for _, v := range s.characterDict {
		if v.EntityId() == int(sync.Entity.Id) {
			v.SetEntityData(sync.GetEntity())
		} else {
			response := &pt.SpaceEntitySyncResponse{EntitySync: sync}
			v.Conn.Send(response)
		}
	}
}

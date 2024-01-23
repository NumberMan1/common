package model

import (
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/vector3"
)

// Entity 在MMO世界进行同步的实体
type Entity struct {
	id        int
	SpaceId   int             //所在地图ID
	Position  vector3.Vector3 //位置
	Direction vector3.Vector3 //方向
}

func NewEntity(id int, position, direction vector3.Vector3) *Entity {
	return &Entity{id: id, Position: position, Direction: direction}
}

func (e *Entity) EntityId() int {
	return e.id
}

func (e *Entity) GetData() *proto.NEntity {
	p := &proto.NEntity{
		Id: int32(e.id),
		Position: &proto.NVector3{
			X: int32(e.Position.X),
			Y: int32(e.Position.Y),
			Z: int32(e.Position.Z),
		},
		Direction: &proto.NVector3{
			X: int32(e.Direction.X),
			Y: int32(e.Direction.Y),
			Z: int32(e.Direction.Z),
		},
	}
	return p
}

func (e *Entity) SetEntityData(entity *proto.NEntity) {
	e.Position.X = float64(entity.Position.X)
	e.Position.Y = float64(entity.Position.Y)
	e.Position.Z = float64(entity.Position.Z)
	e.Direction.X = float64(entity.Direction.X)
	e.Direction.Y = float64(entity.Direction.Y)
	e.Direction.Z = float64(entity.Direction.Z)
}

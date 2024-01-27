package service

import (
	"github.com/NumberMan1/common/game_server/mgr"
	"github.com/NumberMan1/common/game_server/model"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/vector3"
	"math/rand"
	"time"
)

var (
	singleUserService = singleton.Singleton{}
)

type UserService struct {
}

func GetUserServiceInstance() *UserService {
	instance, _ := singleton.GetOrDo[*UserService](&singleUserService, func() (*UserService, error) {
		return &UserService{}, nil
	})
	return instance
}

func (us *UserService) Start() {
	network.GetMessageRouterInstance().Subscribe("proto.GameEnterRequest", network.MessageHandler{Op: us.gameEnterRequest})
}

func (us *UserService) gameEnterRequest(msg network.Msg) {
	logger.SLCInfo("有玩家进入游戏")
	entityId := mgr.GetEntityManagerInstance().NewEntityId()
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	pos := vector3.Vector3{
		X: (r.Float64() * 500) * float64(r.Int()%2),
		Y: 0,
		Z: (r.Float64() * 500) * float64(r.Int()%2),
	}
	pos.Multiply(1000)
	character := model.NewCharacter(entityId, pos, vector3.Zero3())
	//通知玩家登录成功
	response := &proto.GameEnterResponse{
		Success: true,
		Entity:  character.GetData(),
	}
	msg.Sender.Send(response)
	//将新角色加入到地图
	space := GetSpaceServiceInstance().GetSpace(6) //新手村
	space.CharacterJoin(msg.Sender, character)
}

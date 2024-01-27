package service

import (
	"github.com/NumberMan1/common/game_server/model"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

var (
	singleSpaceService = singleton.Singleton{}
)

type SpaceService struct {
	spaceDict map[int]*model.Space
}

func GetSpaceServiceInstance() *SpaceService {
	instance, _ := singleton.GetOrDo[*SpaceService](&singleSpaceService, func() (*SpaceService, error) {
		return &SpaceService{
			spaceDict: map[int]*model.Space{},
		}, nil
	})
	return instance
}

func (ss *SpaceService) Start() {
	//位置同步请求
	network.GetMessageRouterInstance().Subscribe("proto.SpaceEntitySyncRequest", network.MessageHandler{Op: ss.spaceEntitySyncRequest})
	//新手村场景对象
	sp := model.NewSpace(6, "新手村")
	ss.spaceDict[sp.Id] = sp
}

func (ss *SpaceService) GetSpace(id int) *model.Space {
	return ss.spaceDict[id]
}

func (ss *SpaceService) spaceEntitySyncRequest(msg network.Msg) {
	sp := msg.Sender.Get("Space")
	if sp == nil {
		return
	}
	sp.(*model.Space).UpdateEntity(msg.Message.(*proto.SpaceEntitySyncRequest).EntitySync)
}

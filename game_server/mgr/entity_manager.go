package mgr

import (
	"github.com/NumberMan1/common/game_server/model"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/vector3"
	"sync"
)

var (
	singleEntityManager = singleton.Singleton{}
)

// EntityManager Entity管理器
type EntityManager struct {
	index       int
	allEntities map[int]*model.Entity
	mutex       sync.Mutex
}

func GetEntityManagerInstance() *EntityManager {
	instance, _ := singleton.GetOrDo[*EntityManager](&singleEntityManager, func() (*EntityManager, error) {
		return &EntityManager{
			index:       1,
			allEntities: map[int]*model.Entity{},
			mutex:       sync.Mutex{},
		}, nil
	})
	return instance
}

func (em *EntityManager) CreateEntity() *model.Entity {
	em.mutex.Lock()
	entity := &model.Entity{
		SpaceId:   em.index,
		Position:  vector3.Zero3(),
		Direction: vector3.Zero3(),
	}
	em.index += 1
	em.allEntities[entity.EntityId()] = entity
	em.mutex.Unlock()
	return entity
}

func (em *EntityManager) NewEntityId() int {
	em.mutex.Lock()
	id := em.index
	em.index += 1
	em.mutex.Unlock()
	return id
}

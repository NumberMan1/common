package rank

import (
	"sync"
)

type Rank struct {
	ID        uint32
	Category  int
	Cache     sync.Map
	count     int
	blackList []any
}

func (r *Rank) GetMember(key any) any {
	value, ok := r.Cache.Load(key)
	if ok {
		return value
	}
	return nil
}

func (r *Rank) AddMember(key, value any) {
	r.Cache.Store(key, value)
	r.count++
}

func (r *Rank) DelMember(key any) {
	r.Cache.Delete(key)
	r.count--
}

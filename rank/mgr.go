package rank

type Manager struct {
	ranks map[uint32]*Rank
}

func NewManager() *Manager {
	return &Manager{make(map[uint32]*Rank)}
}

func (m *Manager) AddRank(Id uint32, rank *Rank) {
	if _, exist := m.ranks[Id]; exist {
		return
	}
	m.ranks[Id] = rank
}

func (m *Manager) DelRank(Id uint32) {
	delete(m.ranks, Id)
}

func (m *Manager) GetRank(Id uint32) *Rank {
	return m.ranks[Id]
}

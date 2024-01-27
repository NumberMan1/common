package core

type TypeAttributeStore struct {
	dict map[string]any
}

func NewTypeAttributeStore() *TypeAttributeStore {
	return &TypeAttributeStore{dict: map[string]any{}}
}

func (s *TypeAttributeStore) Set(key string, value any) {
	s.dict[key] = value
}

func (s *TypeAttributeStore) Get(key string) any {
	return s.dict[key]
}

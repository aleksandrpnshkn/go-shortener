package store

type MemoryStore struct {
	cache map[string]string
}

func (m *MemoryStore) Get(key string) (value string, isFound bool) {
	value, ok := m.cache[key]
	return value, ok
}

func (m *MemoryStore) Set(key string, value string) error {
	m.cache[key] = value
	return nil
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		cache: map[string]string{},
	}
}

package store

type MemoryStorage struct {
	cache map[string]string
}

func (m *MemoryStorage) Get(key string) (value string, isFound bool) {
	value, ok := m.cache[key]
	return value, ok
}

func (m *MemoryStorage) Set(key string, value string) error {
	m.cache[key] = value
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		cache: map[string]string{},
	}
}

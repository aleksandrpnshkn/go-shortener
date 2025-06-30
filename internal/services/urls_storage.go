package services

type FullURL string

type FullURLsStorage struct {
	storage map[string]string
}

func (s *FullURLsStorage) Set(code Code, url FullURL) {
	s.storage[string(code)] = string(url)
}

func (s *FullURLsStorage) Get(code Code) (FullURL, bool) {
	url, codeExists := s.storage[string(code)]

	return FullURL(url), codeExists
}

func NewFullURLsStorage() *FullURLsStorage {
	fullURLsStorage := FullURLsStorage{
		storage: map[string]string{},
	}

	return &fullURLsStorage
}

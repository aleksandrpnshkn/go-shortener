package services

type Shortener struct {
	codeGenerator CodeGenerator
	urlsStorage   *URLsStorage
	baseURL       string
}

func (s *Shortener) Shorten(URL OriginalURL) string {
	var code Code
	for codeExists := true; codeExists; {
		code = s.codeGenerator.Generate()
		_, codeExists = s.urlsStorage.Get(code)
	}

	s.urlsStorage.Set(code, URL)

	shortURL := s.baseURL + "/" + string(code)

	return shortURL
}

func NewShortener(
	codeGenerator CodeGenerator,
	urlsStorage *URLsStorage,
	baseURL string,
) *Shortener {
	shortener := Shortener{
		codeGenerator,
		urlsStorage,
		baseURL,
	}

	return &shortener
}

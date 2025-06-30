package services

type Shortener struct {
	codeGenerator CodeGenerator
	urlsStorage   FullURLsStorage
	baseURL       string
}

func (s *Shortener) Shorten(URL FullURL) string {
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
	urlsStorage FullURLsStorage,
	baseURL string,
) *Shortener {
	shortener := Shortener{
		codeGenerator,
		urlsStorage,
		baseURL,
	}

	return &shortener
}

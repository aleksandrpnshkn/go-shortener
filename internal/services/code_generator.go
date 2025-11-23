package services

import (
	"math/rand"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

var codeRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

// https://stackoverflow.com/a/31832326
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = codeRunes[rand.Intn(len(codeRunes))]
	}
	return string(b)
}

// CodeGenerator - интерфейс для генерации кодов.
type CodeGenerator interface {
	Generate() types.Code
}

// RandomCodeGenerator генерирует случайные коды.
type RandomCodeGenerator struct {
	codeLength int
}

// Generate генерирурет код.
func (c *RandomCodeGenerator) Generate() types.Code {
	return types.Code(randStringRunes(c.codeLength))
}

// NewRandomCodeGenerator - создаёт новый CodeGenerator.
func NewRandomCodeGenerator(codeLength int) CodeGenerator {
	generator := RandomCodeGenerator{
		codeLength: codeLength,
	}
	return &generator
}

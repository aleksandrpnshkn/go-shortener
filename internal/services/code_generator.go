package services

import "math/rand"

var codeRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

// https://stackoverflow.com/a/31832326
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = codeRunes[rand.Intn(len(codeRunes))]
	}
	return string(b)
}

type Code string

type CodeGenerator struct {
	codeLength int
}

func (c *CodeGenerator) Generate() Code {
	return Code(randStringRunes(c.codeLength))
}

func NewCodeGenerator(codeLength int) *CodeGenerator {
	generator := CodeGenerator{
		codeLength: codeLength,
	}
	return &generator
}

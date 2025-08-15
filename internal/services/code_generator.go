package services

import (
	"math/rand"
	"strconv"
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

type Code string

type CodeGenerator interface {
	Generate() Code
}

type RandomCodeGenerator struct {
	codeLength int
}

func (c *RandomCodeGenerator) Generate() Code {
	return Code(randStringRunes(c.codeLength))
}

func NewRandomCodeGenerator(codeLength int) CodeGenerator {
	generator := RandomCodeGenerator{
		codeLength: codeLength,
	}
	return &generator
}

type TestCodeGenerator struct {
	testCodePrefix string
	counter        int
}

func (c *TestCodeGenerator) Generate() Code {
	c.counter++
	return Code(c.testCodePrefix + strconv.Itoa(c.counter))
}

func NewTestGenerator(testCodePrefix string) CodeGenerator {
	generator := TestCodeGenerator{
		testCodePrefix: testCodePrefix,
		counter:        0,
	}
	return &generator
}

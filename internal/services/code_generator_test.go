package services

import (
	"testing"
)

func BenchmarkRandomCodeGenerator(b *testing.B) {
	codeGenerator := NewRandomCodeGenerator(10)

	for i := 0; i < b.N; i++ {
		codeGenerator.Generate()
	}
}

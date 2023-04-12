package math

var X string = "Obrigado por usado modulo Math"

type math struct {
	A int
	B int
}

func NewMath(a, b int) math {
	return math{A: a, B: b}
}

func (m math) Add() int {
	return m.A + m.B
}

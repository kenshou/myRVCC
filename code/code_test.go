package code

import (
	"myRVCC/lexer"
	"myRVCC/parser"
	"testing"
)

func TestGenCode(t *testing.T) {
	l := lexer.New("{ a=3; z=5;return a+z; }")
	p := parser.New(l)
	program := p.ParseProgram()
	genCode(program)
}

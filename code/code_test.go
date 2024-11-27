package code

import (
	"myRVCC/lexer"
	"myRVCC/parser"
	"testing"
)

func TestGenCode(t *testing.T) {
	l := lexer.New("{ x=3; y=5; *(&y-8)=7; return x; }")
	p := parser.New(l)
	program := p.ParseProgram()
	genCode(program)
}

package code

import (
	"myRVCC/lexer"
	"myRVCC/parser"
	"testing"
)

func TestGenCode(t *testing.T) {
	l := lexer.New("{ i=0; while(i<10) { i=i+1; } return i; }")
	p := parser.New(l)
	program := p.ParseProgram()
	genCode(program)
}

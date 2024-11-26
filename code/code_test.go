package code

import (
	"myRVCC/lexer"
	"myRVCC/parser"
	"testing"
)

func TestGenCode(t *testing.T) {
	l := lexer.New("a=5")
	p := parser.New(l)
	program := p.ParseProgram()
	genCode(program)
}

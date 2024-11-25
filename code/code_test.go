package code

import (
	"myRVCC/lexer"
	"myRVCC/parser"
	"testing"
)

func TestGenCode(t *testing.T) {
	l := lexer.New("1+2+(3+4)*5")
	p := parser.New(l)
	program := p.ParseProgram()
	genCode(program)
}

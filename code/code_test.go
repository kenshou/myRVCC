package code

import (
	"myRVCC/lexer"
	"myRVCC/parser"
	"testing"
)

func TestGenCode(t *testing.T) {
	l := lexer.New("for(a=1;a<10;a=a+1){return 0;}")
	p := parser.New(l)
	program := p.ParseProgram()
	genCode(program)
}

package code

import (
	"myRVCC/lexer"
	"myRVCC/parser"
	"testing"
)

func TestGenCode(t *testing.T) {
	l := lexer.New("if(0){10}else{20}return 0;")
	p := parser.New(l)
	program := p.ParseProgram()
	genCode(program)
}

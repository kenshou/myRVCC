package parser

import (
	"myRVCC/lexer"
	"myRVCC/logger"
	"testing"
)

func TestParserString(t *testing.T) {
	input := "(1 + 2) * 3/4+2/3"
	p := New(lexer.New(input))
	program := p.ParseProgram()
	logger.Info(program.String())
}

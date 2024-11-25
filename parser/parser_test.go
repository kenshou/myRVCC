package parser

import (
	"myRVCC/lexer"
	"myRVCC/logger"
	"testing"
)

func TestParserString(t *testing.T) {
	input := "a=b=10;a=c+10"
	p := New(lexer.New(input))
	program := p.ParseProgram()
	logger.Info(program.String())
}

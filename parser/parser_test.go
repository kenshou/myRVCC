package parser

import (
	"myRVCC/lexer"
	"myRVCC/logger"
	"testing"
)

func TestParserString(t *testing.T) {
	input := "10;"
	p := New(lexer.New(input))
	program := p.ParseProgram()
	logger.Info(program.String())
}

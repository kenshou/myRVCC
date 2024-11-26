package parser

import (
	"myRVCC/ast"
	"myRVCC/lexer"
	"myRVCC/logger"
	"myRVCC/token"
	"testing"
)

func TestParserString(t *testing.T) {
	input := "a=b=10;a=c+10"
	p := New(lexer.New(input))
	program := p.ParseProgram()
	logger.Info(program.String())
}

func TestEnv(t *testing.T) {
	env := ast.CreateEnv()
	ident := &token.Token{
		Kind:    token.IDENT,
		Literal: "a",
	}
	obj := env.FindOrCreateIdentifier(ident)
	logger.Info("obj: %+v", obj)
	obj = env.FindOrCreateIdentifier(ident)
	logger.Info("obj: %+v", obj)
}

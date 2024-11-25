package main

import (
	"myRVCC/code"
	"myRVCC/lexer"
	"myRVCC/logger"
	"myRVCC/parser"
	"myRVCC/token"
	"os"
	"strconv"
)

func main() {
	//判断传入的参数是否为2个，第一个为程序名称，第二个位传入的参数
	if len(os.Args) < 2 {
		logger.Error("%s: invalid number of arguments\n", os.Args[0])
		return
	}
	//exp为求值的算式
	exp := os.Args[1]
	run(exp)
	return
}

func run(exp string) {
	l := lexer.New(exp)
	p := parser.New(l)
	program := p.ParseProgram()

	code.GenRootCode(program)
}
func getNumber(tok token.Token) int64 {
	if tok.Kind != token.INT {
		logger.Panic("[%s] getNumber: token is not int", tok.Literal)
	}
	value, err := strconv.ParseInt(tok.Literal, 10, 64)
	if err != nil {
		logger.Panic("[%s] getNumber: Atoi error", tok.Literal)
	}
	return value
}

package main

import (
	"myRVCC/asm"
	"myRVCC/lexer"
	"myRVCC/logger"
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

	asm.Globl("main")
	asm.Label("main")
	l := lexer.New(exp)
	//处理第一个数字
	tok := l.NextToken()
	asm.Li(asm.REG_A0, getNumber(tok))
	isNeg := false
l_for:
	for {
		tok = l.NextToken()
		switch tok.Kind {
		case token.EOF:
			break l_for //跳出整个循环
		case token.ADD:
			isNeg = false
		case token.SUB:
			isNeg = true
		case token.INT:
			num := getNumber(tok)
			if isNeg {
				num = -num
			}
			asm.Addi(asm.REG_A0, asm.REG_A0, num)
		}
	}
	asm.Ret()
	return
}
func getNumber(tok token.Token) int {
	if tok.Kind != token.INT {
		logger.Panic("[%s] getNumber: token is not int", tok.Literal)
	}
	value, err := strconv.Atoi(tok.Literal)
	if err != nil {
		logger.Panic("[%s] getNumber: Atoi error", tok.Literal)
	}
	return value
}

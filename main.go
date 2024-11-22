package main

import (
	"fmt"
	"myRVCC/logger"
	"myRVCC/utils"
	"os"
)

func main() {
	//判断传入的参数是否为2个，第一个为程序名称，第二个位传入的参数
	if len(os.Args) < 2 {
		logger.Error("%s: invalid number of arguments\n", os.Args[0])
		return
	}
	//exp为求值的算式
	exp := os.Args[1]

	fmt.Println("	.globl main")
	fmt.Println("main:")
	//构建中间表达式
	num := 0
	isNeg := false
	//是否为第一个初始化的数字
	is1st := true
	for _, p := range exp {
		switch p {
		case '-':
			if num != 0 {
				printNum(num, &is1st)
			}
			isNeg = true
			num = 0
		case '+':
			if num != 0 {
				printNum(num, &is1st)
			}
			isNeg = false
			num = 0
		default:
			nowNum := int(p - '0')
			if isNeg {
				nowNum = -nowNum
			}
			num = num*10 + nowNum
		}
	}
	printNum(num, &is1st)
	fmt.Println("	ret")
	return
}

func printNum(num int, is1st *bool) {
	if *is1st {
		utils.PrintLine("	li a0,%d", num)
		*is1st = false
	} else {
		utils.PrintLine("	addi a0,a0,%d", num)
	}

}

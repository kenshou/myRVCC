package main

import (
	"fmt"
	"myRVCC/logger"
	"os"
)

func main() {
	//判断传入的参数是否为2个，第一个为程序名称，第二个位传入的参数
	if len(os.Args) < 2 {
		logger.Error("%s: invalid number of arguments\n", os.Args[0])
		return
	}
	fmt.Println("	.globl main")
	fmt.Println("main:")
	fmt.Println(fmt.Sprintf("	li a0,%s", os.Args[1]))
	fmt.Println("	ret")
	return
}

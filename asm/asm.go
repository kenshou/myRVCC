package asm

import "myRVCC/utils"

func Globl(name string) {
	utils.PrintLine("	.globl %s", name)
}

func Label(name string) {
	utils.PrintLine("%s:", name)
}

func Li(rd REGISTER, imm int) {
	utils.PrintLine("	li %s,%d", rd, imm)
}

func Addi(rd REGISTER, rs REGISTER, imm int) {
	utils.PrintLine("	addi %s,%s,%d", rd, rs, imm)
}

func Ret() {
	utils.PrintLine("	ret")
}

package asm

import "myRVCC/utils"

func Globl(name string) {
	utils.PrintLine("	.globl %s", name)
}

func Label(name string) {
	utils.PrintLine("%s:", name)
}

func Li(rd REGISTER, imm int64) {
	utils.PrintLine("	li %s,%d", rd, imm)
}

func Addi(rd REGISTER, rs REGISTER, imm int64) {
	utils.PrintLine("	addi %s,%s,%d", rd, rs, imm)
}
func Add(rd REGISTER, rs1 REGISTER, rs2 REGISTER) {
	utils.PrintLine("	add %s,%s,%s", rd, rs1, rs2)
}
func Sub(rd REGISTER, rs1 REGISTER, rs2 REGISTER) {
	utils.PrintLine("	sub %s,%s,%s", rd, rs1, rs2)
}

func Multi(rd REGISTER, rs1 REGISTER, rs2 REGISTER) {
	utils.PrintLine("	mul %s,%s,%s", rd, rs1, rs2)
}

func Div(rd REGISTER, rs1 REGISTER, rs2 REGISTER) {
	utils.PrintLine("	div %s,%s,%s", rd, rs1, rs2)
}

// Neg 对寄存器取反
func Neg(rd REGISTER, rs REGISTER) {
	utils.PrintLine("	neg %s,%s", rd, rs)
}

func Xor(rd REGISTER, rs1 REGISTER, rs2 REGISTER) {
	utils.PrintLine("	xor %s,%s,%s", rd, rs1, rs2)
}
func Xori(rd REGISTER, rs1 REGISTER, imm int64) {
	utils.PrintLine(" xori %s, %s, %d", rd, rs1, imm)
}

// Seqz 等于0则置1
func Seqz(rd REGISTER, rs1 REGISTER) {
	utils.PrintLine("	seqz %s,%s", rd, rs1)
}

// Snez 不等于0则置1
func Snez(rd REGISTER, rs1 REGISTER) {
	utils.PrintLine("	snez %s,%s", rd, rs1)
}

// Slt 如果 rs1<rs2则 rd 置1，否则置0
func Slt(rd REGISTER, rs1 REGISTER, rs2 REGISTER) {
	utils.PrintLine("	slt %s,%s,%s", rd, rs1, rs2)
}
func Ret() {
	utils.PrintLine("	ret")
}

// PushA0 压栈，将结果临时压入栈中备用
// sp为栈指针，栈反向向下增长，64位下，8个字节为一个单位，所以sp-8
// 当前栈指针的地址就是sp，将a0的值压入栈
// 不使用寄存器存储的原因是因为需要存储的值的数量是变化的。
func PushA0() {
	Addi(REG_SP, REG_SP, -8)
	utils.PrintLine("	sd a0,0(sp)")
}

// Pop 弹栈，将sp指向的地址的值，弹出到a1
func Pop(reg REGISTER) {
	utils.PrintLine("	ld %s,0(sp)", reg)
	Addi(REG_SP, REG_SP, 8)
}

package code

import (
	"myRVCC/asm"
	"myRVCC/ast"
	"myRVCC/token"
)

func GenCode(node ast.Node) {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		asm.Li(asm.REG_A0, node.Value)
	case *ast.InfixExpression:
		genCodeInfixExpression(node)
	case *ast.PrefixExpression:
		genCodePrefixExpression(node)
	case *ast.Program:
		genCodeProgram(node)
	case *ast.ExpressionStatement:
		GenCode(node.Expression)
	default:
		panic("unsupported node type")
	}
}

func genCodePrefixExpression(node *ast.PrefixExpression) {
	switch node.Token.Kind {
	case token.SUB:
		GenCode(node.Right)
		asm.Neg(asm.REG_A0, asm.REG_A0)
	case token.ADD:
		GenCode(node.Right)
		return
	default:
		panic("genCodePrefixExpression unsupported operator ")
	}
}

func genCodeProgram(program *ast.Program) {
	for _, statement := range program.Statements {
		GenCode(statement)
	}
}

func genCodeInfixExpression(node *ast.InfixExpression) {
	//先递归右节点存入堆栈，再递归左节点到A0;然后弹出右节点的值到A1。
	GenCode(node.Right)
	asm.PushA0()
	GenCode(node.Left)
	asm.Pop(asm.REG_A1)
	switch node.Token.Kind {
	case token.ADD:
		asm.Add(asm.REG_A0, asm.REG_A0, asm.REG_A1)
	case token.SUB:
		asm.Sub(asm.REG_A0, asm.REG_A0, asm.REG_A1)
	case token.MUL:
		asm.Multi(asm.REG_A0, asm.REG_A0, asm.REG_A1)
	case token.DIV:
		asm.Div(asm.REG_A0, asm.REG_A0, asm.REG_A1)
	case token.EQ:
		asm.Xor(asm.REG_A0, asm.REG_A0, asm.REG_A1)
		asm.Seqz(asm.REG_A0, asm.REG_A0)
	case token.NEQ:
		asm.Xor(asm.REG_A0, asm.REG_A0, asm.REG_A1)
		asm.Snez(asm.REG_A0, asm.REG_A0)
	case token.LT:
		asm.Slt(asm.REG_A0, asm.REG_A0, asm.REG_A1)
	case token.LEQ:
		asm.Slt(asm.REG_A0, asm.REG_A1, asm.REG_A0)
		asm.Xori(asm.REG_A0, asm.REG_A0, 1)
	case token.GT:
		asm.Slt(asm.REG_A0, asm.REG_A1, asm.REG_A0)
	case token.GEQ:
		asm.Slt(asm.REG_A0, asm.REG_A0, asm.REG_A1)
		asm.Xori(asm.REG_A0, asm.REG_A0, 1)
	default:
		panic("unsupported operator")
	}
}

package code

import (
	"myRVCC/asm"
	"myRVCC/ast"
	"myRVCC/logger"
	"myRVCC/token"
)

const (
	ReturnLabel = ".L.return"
)

func GenRootCode(node ast.Node) {
	genCode(node)
}

func genCode(node ast.Node) {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		//asm.Comment("%d", node.Value)
		asm.Li(asm.REG_A0, node.Value)
	case *ast.InfixExpression:
		genCodeInfixExpression(node)
	case *ast.PrefixExpression:
		genCodePrefixExpression(node)
	case *ast.Program:
		genCodeProgram(node)
	case *ast.ExpressionStatement:
		genCode(node.Expression)
	case *ast.Identifier:
		genCodeIdentifier(node)
	case *ast.ReturnStatement:
		genCodeReturnStatement(node)
	case *ast.BlockStatement:
		genCodeBlockStatement(node)

	default:
		panic("unsupported node type")
	}
}

func genCodeBlockStatement(node *ast.BlockStatement) {
	for _, stmt := range node.Statements {
		genCode(stmt)
	}
}

func genCodeReturnStatement(node *ast.ReturnStatement) {
	genCode(node.ReturnValue)
	asm.J(ReturnLabel)
}

func genCodeIdentifier(node *ast.Identifier) {
	// 计算出变量的地址，然后存入a0
	genAddress(node)
	// 访问a0地址中存储的数据，存入到a0当中
	asm.Ld(asm.REG_A0, 0, asm.REG_A0)
}

func genCodePrefixExpression(node *ast.PrefixExpression) {
	switch node.Token.Kind {
	case token.SUB:
		genCode(node.Right)
		asm.Neg(asm.REG_A0, asm.REG_A0)
	case token.ADD:
		genCode(node.Right)
		return
	default:
		panic("genCodePrefixExpression unsupported operator ")
	}
}

func genCodeProgram(program *ast.Program) {
	env := program.Env
	AssignVarOffset(env)
	asm.Globl("main")
	asm.Label("main")
	// 栈布局
	//-------------------------------// sp
	//              fp
	//-------------------------------// fp = sp-8
	//             变量
	//-------------------------------// sp = sp-8-StackSize
	//           表达式计算
	//-------------------------------//

	// 将fp压入栈中，保存fp的值
	asm.Addi(asm.REG_SP, asm.REG_SP, -8)
	asm.Sd(asm.REG_FP, 0, asm.REG_SP)
	//将sp写入fp
	asm.Mv(asm.REG_FP, asm.REG_SP)
	//26个字母*8字节=208字节，栈腾出208字节的空间
	asm.Addi(asm.REG_SP, asm.REG_SP, -env.StackSize)
	for _, statement := range program.Statements {
		genCode(statement)
	}
	asm.Label(ReturnLabel)
	//恢复栈sp
	asm.Mv(asm.REG_SP, asm.REG_FP)
	//将最早的fp保存的值弹栈，恢复fp的值
	asm.Ld(asm.REG_FP, 0, asm.REG_SP)
	asm.Addi(asm.REG_SP, asm.REG_SP, 8)
	//返回
	asm.Ret()
}

func AssignVarOffset(env *ast.Env) {
	offset := int64(0)
	for _, obj := range env.IdentObjArr {
		// 每个变量分配8字节
		offset += 8
		// 为每个变量赋一个偏移量，或者说是栈中地址
		obj.Offset = -offset
	}
	// 将栈对齐到16字节
	env.StackSize = AlignTo(offset, 16)
}
func AlignTo(n, align int64) int64 {
	return (n + align - 1) / align * align
}
func genCodeInfixExpression(node *ast.InfixExpression) {
	asm.Comment("%s %s %s", node.Left.TokenLiteral(), node.Operator, node.Right.TokenLiteral())
	//先递归右节点存入堆栈，再递归左节点到A0;然后弹出右节点的值到A1。
	genCode(node.Right)
	asm.PushA0()
	genCode(node.Left)
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
	case token.ASSIGN:
		genCodeAssign(node)
	default:
		panic("unsupported operator")
	}
}

func genCodeAssign(node *ast.InfixExpression) {
	if node.Token.Kind == token.ASSIGN {
		left, ok := node.Left.(*ast.Identifier)
		if ok {
			genAddress(left)
			asm.PushA0()
			genCode(node.Right)
			asm.Pop(asm.REG_A1)
			asm.Sd(asm.REG_A0, 0, asm.REG_A1)
			return
		} else {
			logger.Panic("genCodeAssign:unsupported operator %+v", node)
		}
	} else {
		logger.Panic("genCodeAssign:unsupported operator %+v", node)
	}
}

func genAddress(identifier *ast.Identifier) {
	if identifier.Token.Kind == token.IDENT {
		offset := identifier.Obj.Offset
		asm.Addi(asm.REG_A0, asm.REG_FP, offset)
	} else {
		logger.Panic("unsupported operator %+v", identifier)
	}
}

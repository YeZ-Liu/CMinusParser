// Copyright 2020. All rights reserved.
// Author: Zhifei Liu, 2020/6
// Filename: scanner.go
// Package: scan
// Description: 本文件定义了辅助打印函数,包括词法分析和语法分析
package scan

import (
	"fmt"
	"os"
)

// 将词法分析结果输入到文件中
func HelpPrintFile(token Token, lexeme TokenString, line int, file *os.File) {
	fmt.Fprintf(file, "[Line %d]:", line)
	switch token {
	case ID:
		fmt.Fprint(file, " ID ---->  ")
	case NUM:
		fmt.Fprint(file, " NUM ---->  ")
	case IF, ELSE, WHILE, INT, VOID, RETURN:
		fmt.Fprint(file, " KEY ---->  ")
	case PLUS, MINUS, MUL, DIV, LT, GT, LE, GE, EQ, ASSIGN, NOT_EQ:
		fmt.Fprint(file, " OP ---->  ")
	case SEMI, COMMA, L_PARE_L, L_PARE_M, L_PARE_S, R_PARE_L, R_PARE_M, R_PARE_S:
		fmt.Fprint(file, " SEP ---->  ")
	case COMMENT:
		fmt.Fprint(file, " COMMENT ---->  ")
	case ERROR:
		fmt.Fprint(file, " ERROR ---->  ")
	case EOF_TOKEN:
		fmt.Fprint(file, " EOF ")
	}
	fmt.Fprintln(file, string(lexeme))
}

// 打印抽象语法树
func HelpPrintTree(root *ASTNode, sp int, c byte, file *os.File) {
	if root == nil {
		return
	}

	for i := 0; i < 2*sp; i++ {
		fmt.Fprintf(file, "%c", c)
	}
	HelpPrintNode(root, file)
	fmt.Fprintln(file)
	fmt.Fprintln(file)
	if root.left != nil {
		HelpPrintTree(root.left, sp+2, c, file)
	}
	if root.mid != nil {
		HelpPrintTree(root.mid, sp+2, c, file)
	}
	if root.right != nil {
		HelpPrintTree(root.right, sp+2, c, file)
	}
	if root.sibling != nil {
		HelpPrintTree(root.sibling, sp, c, file)
	}
}

// 打印语法树节点信息
func HelpPrintNode(root *ASTNode, file *os.File) {
	fmt.Fprintf(file, "Line: %d  ", root.line)
	switch root.nodeK {
	case STATEMENT: // 语句
		switch root.nodeT {
		case VAR_DECLARATION:
			fmt.Fprint(file, "([VAR_DECLARATION] ")
			fmt.Fprintf(file, "ID:%s; ", string(root.attribute.(TokenString)))
			//switch root.varT {
			//case VAR_TYPE_VOID:
			//	fmt.Fprint(file,"Type: void; ")
			//case VAR_TYPE_INT:
			//	fmt.Fprint(file,"Type: int; ")
			//case VAR_TYPE_INT_VECTOR:
			//	fmt.Fprint(file,"Type: int[]; ")
			//}
			fmt.Fprint(file, ")")
		case FUNC_DECLARATION:
			fmt.Fprint(file, "([FUNC_DECLARATION]  ")
			fmt.Fprintf(file, "ID:%s; ", string(root.attribute.(TokenString)))
			//switch root.varT {
			//case VAR_TYPE_VOID:
			//	fmt.Fprint(file,"Return type: void")
			//case VAR_TYPE_INT:
			//	fmt.Fprint(file,"Return type: int")
			//case VAR_TYPE_INT_VECTOR:
			//	fmt.Fprint(file,"Return type: int[]")
			//}
			fmt.Fprint(file, ")")
		case COMPOUND:
			fmt.Fprint(file, "([COMPOUND] ")
			fmt.Fprint(file, ")")
		case SELECTION_STMT:
			fmt.Fprint(file, "([SELECTION] ")
			fmt.Fprint(file, ")")
		case ITERATION_STMT:
			fmt.Fprint(file, "([ITERATION] ")
			fmt.Fprint(file, ")")
		case RETURN_STMT:
			fmt.Fprint(file, "([RETURN] ")
			fmt.Fprint(file, ")")
		}
	case EXPRESSION: // 表达式
		switch root.nodeT {
		case VAR:
			fmt.Fprint(file, "([VAR] ")
			fmt.Fprintf(file, "ID:%s; ", string(root.attribute.(TokenString)))
			fmt.Fprint(file, ")")
		case ASSIGNMENT:
			fmt.Fprint(file, "([ASSIGNMENT] ")
			fmt.Fprint(file, ")")
		case CALL:
			fmt.Fprint(file, "([CALL] ")
			fmt.Fprintf(file, "ID:%s; ", string(root.attribute.(TokenString)))
			if root.left == nil {
				fmt.Fprintf(file, "Args: None")
			}
			fmt.Fprint(file, ")")
		case COMPARE:
			fmt.Fprint(file, "([COMPARE] ")
			HelpPrintToken(root.attribute.(Token), file)
			fmt.Fprint(file, ")")
		case CONST:
			fmt.Fprint(file, "([CONST] ")
			fmt.Fprintf(file, "VALUE:%d; ", root.attribute.(int64))
			fmt.Fprint(file, ")")
		case OPERATION:
			fmt.Fprint(file, "([OPERATION] ")
			HelpPrintToken(root.attribute.(Token), file)
			fmt.Fprint(file, ")")
		}
	case PARAMS: // 形式参数
		fmt.Fprint(file, "([PARAMS] ")
		if root.left == nil {
			fmt.Fprintf(file, "void")
		}
		fmt.Fprint(file, ")")
	case ARGS: // 实参
		fmt.Fprint(file, "([ARGS] ")
		if root.left == nil {
			fmt.Fprintf(file, "none")
		}
		fmt.Fprint(file, ")")
	case PARAM:
		fmt.Fprint(file, "([PARAM] ")
		fmt.Fprintf(file, "ID:%s; ", string(root.attribute.(TokenString)))
		//switch root.varT {
		//case VAR_TYPE_VOID:
		//	fmt.Fprint(file,"Type: void")
		//case VAR_TYPE_INT:
		//	fmt.Fprint(file,"Type: int")
		//case VAR_TYPE_INT_VECTOR:
		//	fmt.Fprint(file,"Type: int[]")
		//}
		fmt.Fprint(file, ")")
	case TYPE:
		fmt.Fprint(file, "([TYPE] ")
		switch root.varT {
		case VAR_TYPE_VOID:
			fmt.Fprint(file, "Type: void")
		case VAR_TYPE_INT:
			fmt.Fprint(file, "Type: int")
		case VAR_TYPE_INT_VECTOR:
			fmt.Fprint(file, "Type: int[]")
		}
		fmt.Fprint(file, ")")
	}
}

// 根据类型打印token对应的符号
func HelpPrintToken(t Token, file *os.File) {
	switch t {
	case INT:
		fmt.Fprint(file, "int")
	case VOID:
		fmt.Fprint(file, "void")
	case MINUS:
		fmt.Fprint(file, "-")
	case PLUS:
		fmt.Fprint(file, "+")
	case MUL:
		fmt.Fprint(file, "*")
	case DIV:
		fmt.Fprint(file, "/")
	case LT:
		fmt.Fprint(file, "<")

	case LE:
		fmt.Fprint(file, "<=")

	case GT:
		fmt.Fprint(file, ">")

	case GE:
		fmt.Fprint(file, ">=")

	case EQ:
		fmt.Fprint(file, "==")

	case NOT_EQ:
		fmt.Fprint(file, "!=")
	}
}

// 打印符号表信息
func HelpPrintTable(root *SymbolTableNode, sp int, c byte, file *os.File) {
	for key, value := range root.data {
		val := value.(*content)
		for i := 0; i < sp; i++ {
			fmt.Fprintf(file, "%c", c)
		}
		fmt.Fprintf(file, "line: %d ", val.line)
		fmt.Fprint(file, "(", key, ")")
	}

	if root.next != nil {
		HelpPrintTable(root.next, sp+4, c, file)
	}

	if root.rightSib != nil {
		HelpPrintTable(root.rightSib, sp, c, file)
	}
}

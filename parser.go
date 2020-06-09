// Copyright 2020. All rights reserved.
// Author: Zhifei Liu, 2020/6
// Filename: parser.go
// Package: scan
// Description: 本文件定义了语法分析器类以及成员函数和初始化工厂函数

package scan

import (
	"fmt"
	"strconv"
)

/**
Parser类定义
*/
type Parser struct {
	buffer     *Buffer     // 输入缓冲区
	aheadToken Token       // 前向Token
	lexeme     TokenString // 扫描出的词素
	scanner    *Scanner    // 词法分析器
}

// 语法分析器工厂函数
func NewParser(buffer *Buffer) *Parser {
	var parser Parser
	parser.buffer = buffer
	parser.scanner = NewScanner(buffer)
	return &parser
}

// 分析类型节点
func (parser *Parser) typeSpecifier() *ASTNode {
	var node *ASTNode
	if parser.aheadToken == INT || parser.aheadToken == VOID {
		node = NewASTNode(TYPE, nil, parser.buffer.Lines())
		node.SetType(parser.aheadToken)
		parser.match(parser.aheadToken)
	}
	return node
}

// 分析声明语句序列
func (parser *Parser) declarationList() *ASTNode {
	var res, cur, next *ASTNode
	cur = parser.declaration()
	res = cur
	// 如果还有另一个声明,则需要继续识别
	for parser.aheadToken == INT || parser.aheadToken == VOID {
		next = parser.declaration()
		cur.SetSibling(next)
		cur = next
	}
	return res
}

// 分析语句，有变量声明语句和函数声明语句
func (parser *Parser) declaration() *ASTNode {
	var identifier TokenString // ID 对应的词素
	var node, typeNode *ASTNode
	line := parser.buffer.Lines() // 当前行数

	typeNode = parser.typeSpecifier() // 类型节点
	identifier = parser.lexeme
	AddIdentifier(string(identifier)) // 添加到符号表
	parser.match(ID)
	// 根据后一个token类型判断是函数声明还是变量声明，以及是否数组声明
	// 函数声明
	if parser.aheadToken == L_PARE_S {
		MoveDown()                                           // 函数参数也属于下一层作用域
		node = NewASTNode(STATEMENT, FUNC_DECLARATION, line) // 表达式类型、函数声明
		parser.match(L_PARE_S)
		p := parser.params()
		parser.match(R_PARE_S)
		c := parser.compoundStmt()
		MoveUp() // 作用域离开函数声明
		node.SetMid(p)
		node.SetRight(c)
	} else { // 变量、数组声明
		node = NewASTNode(STATEMENT, VAR_DECLARATION, line) // 表达式类型、变量声明
		if parser.aheadToken == L_PARE_M {                  // 为数组变量声明
			parser.match(L_PARE_M)
			num := parser.factor() // 将数组大小作为子节点返回
			node.SetRight(num)
			typeNode.SetVec()
			parser.match(R_PARE_M)
		}
		parser.match(SEMI) // 末尾的分号
	}
	node.SetAttr(identifier) // 设置属性为ID
	node.SetLeft(typeNode)   // 设置类型节点

	//}
	return node
}

// 函数形参列表,可以为空VOID
// 按照语法规则来看,需要向前看两个token才能确定是VOID无参数，还是有多个参数
func (parser *Parser) params() *ASTNode {
	var res, cur, next, typeNode *ASTNode
	line := parser.buffer.Lines()

	typeNode = parser.typeSpecifier()
	// 说明是有参数的函数
	if parser.aheadToken == ID {
		AddIdentifier(string(parser.lexeme)) // 参数标识符添加到符号表
		cur = NewASTNode(PARAM, nil, line)   // 单个参数
		cur.SetAttr(parser.lexeme)           // 设置参数ID
		cur.SetLeft(typeNode)                // 设置类型
		parser.match(ID)

		// 判断是否是数组参数
		if parser.aheadToken == L_PARE_M {
			parser.match(L_PARE_M)
			parser.match(R_PARE_M)
			typeNode.SetVec() // 数组参数类型
		}
		res = cur // 第一个参数

		// 其余并列的形式参数
		for parser.aheadToken == COMMA {
			parser.match(COMMA)
			next = parser.param()
			cur.SetSibling(next)
			cur = next
		}
	}

	// 若为没有参数的函数,则参数部分为nil，否则为并列的单个参数
	tmp := NewASTNode(PARAMS, nil, line)
	tmp.SetLeft(res)
	return tmp
}

// 形式参数
func (parser *Parser) param() *ASTNode {
	var cur, typeNode *ASTNode
	line := parser.buffer.Lines()

	typeNode = parser.typeSpecifier()
	cur = NewASTNode(PARAM, nil, line)
	cur.SetLeft(typeNode)

	cur.SetAttr(parser.lexeme)           // 设置形参ID
	AddIdentifier(string(parser.lexeme)) // 参数标识符添加到符号表
	parser.match(ID)

	// 判断是否是数组参数
	if parser.aheadToken == L_PARE_M {
		parser.match(L_PARE_M)
		parser.match(R_PARE_M)
		typeNode.SetVec() // 数组参数
	}
	//}
	return cur
}

// 复合语句
func (parser *Parser) compoundStmt() *ASTNode {
	var cur *ASTNode
	line := parser.buffer.Lines()

	parser.match(L_PARE_L)
	// 中间的局部变量声明和语句序列都是可选的
	// 局部变量声明
	cur = NewASTNode(STATEMENT, COMPOUND, line) // 语句、复合语句

	switch parser.aheadToken {
	case INT, VOID: // 局部变量域
		t := parser.localDeclarations()
		cur.SetLeft(t)
		fallthrough // 两种情况都要判断
	case SEMI, ID, L_PARE_S, NUM, L_PARE_L, IF, WHILE, RETURN:
		// 声明序列
		t := parser.statementList()
		cur.SetRight(t)
	}
	parser.match(R_PARE_L)

	return cur
}

// 声明语句序列
func (parser *Parser) localDeclarations() *ASTNode {
	var res, cur, next *ASTNode
	if parser.aheadToken == INT || parser.aheadToken == VOID {
		cur = parser.varDeclaration()
		res = cur
	} else {
		return res
	}

	for parser.aheadToken == INT || parser.aheadToken == VOID {
		next = parser.varDeclaration()
		cur.SetSibling(next)
		cur = next
	}
	return res
}

// 变量声明
func (parser *Parser) varDeclaration() *ASTNode {
	var node, typeNode *ASTNode
	line := parser.buffer.Lines()

	node = NewASTNode(STATEMENT, VAR_DECLARATION, line) // 语句、变量声明语句
	typeNode = parser.typeSpecifier()
	node.SetLeft(typeNode)

	node.SetAttr(parser.lexeme)          // 设置变量ID属性
	AddIdentifier(string(parser.lexeme)) // 变量标识符添加到符号表
	parser.match(ID)

	if parser.aheadToken == L_PARE_M { // 数组声明
		parser.match(L_PARE_M)
		num := parser.factor()
		node.SetRight(num) // 数组大小
		typeNode.SetVec()
		parser.match(R_PARE_M)
	}
	parser.match(SEMI)
	return node
}

// 语句序列
func (parser *Parser) statementList() *ASTNode {
	var res, cur, next *ASTNode
	switch parser.aheadToken {
	case SEMI, ID, L_PARE_S, NUM, L_PARE_L, IF, WHILE, RETURN:
		cur = parser.statement()
		res = cur
	L1:
		for {
			switch parser.aheadToken {
			case SEMI, ID, L_PARE_S, NUM, L_PARE_L, IF, WHILE, RETURN:
				next = parser.statement()
				cur.SetSibling(next)
				cur = next
			default:
				break L1
			}
		}
	}
	return res
}

// 语句
func (parser *Parser) statement() *ASTNode {
	var res *ASTNode
	switch parser.aheadToken {
	case SEMI, ID, L_PARE_S, NUM: // 表达式语句
		res = parser.expressionStmt()
	case L_PARE_L: // 复合语句
		res = parser.compoundStmt()
	case IF: // 选择语句
		res = parser.selectionStmt()
	case WHILE: // 循环语句
		res = parser.iterationStmt()
	case RETURN: // 返回语句
		res = parser.returnStmt()
	default:
		parser.syntaxError()
	}
	return res
}

// 表达式语句
func (parser *Parser) expressionStmt() *ASTNode {
	var res *ASTNode
	switch parser.aheadToken {
	case SEMI:
		parser.match(SEMI)
	case ID, NUM, L_PARE_S:
		res = parser.expression()
		parser.match(SEMI)
	default:
		parser.syntaxError()
	}
	return res
}

// 表达式语句分为赋值语句和简单算术表达式
// 需要先判断是否是var,然后再通过后面是否有"="判断到底是赋值还是简单算术表达式
func (parser *Parser) expression() *ASTNode {
	var res, cur, next *ASTNode
	var id TokenString
	line := parser.buffer.Lines()

	switch parser.aheadToken {
	case L_PARE_S, NUM: // 可以确定是factor
		cur = parser.simpleExpression()
		res = cur
	case ID: // 有三种可能，一是var,二是call,三是赋值语句，需要后面再进行判断
		id = parser.lexeme
		AddIdentifier(string(parser.lexeme)) // 标识符添加到符号表
		parser.match(ID)

		switch parser.aheadToken {
		case L_PARE_S: // 函数调用
			cur = NewASTNode(EXPRESSION, CALL, line) // 表达式、函数调用
			cur.SetAttr(id)                          // 设置函数ID属性
			parser.match(L_PARE_S)
			t := parser.args()
			cur.SetLeft(t)
			parser.match(R_PARE_S)

			// cur 作为factor 继续向上分析
			for parser.aheadToken == MUL || parser.aheadToken == DIV {
				next = NewASTNode(EXPRESSION, OPERATION, line) // 表达式，操作符
				next.SetLeft(cur)
				next.SetAttr(parser.aheadToken) // 设置操作符属性
				parser.match(parser.aheadToken) // 匹配操作符

				right := parser.factor()
				next.SetRight(right)
				cur = next
			}
			// cur 作为term 继续向上分析
			for parser.aheadToken == PLUS || parser.aheadToken == MINUS {
				next = NewASTNode(EXPRESSION, OPERATION, line) // 表达式，操作符
				next.SetLeft(cur)
				next.SetAttr(parser.aheadToken) // 设置操作符属性
				parser.match(parser.aheadToken) // 匹配操作符

				right := parser.term()
				next.SetRight(right)
				cur = next
			}
			// cur 作为additive_expression 继续向上分析
			if parser.aheadToken == LE || parser.aheadToken == LT || parser.aheadToken == GT || parser.aheadToken == GE || parser.aheadToken == EQ || parser.aheadToken == NOT_EQ {
				next = NewASTNode(EXPRESSION, COMPARE, line) // 表达式，比较语句
				next.SetAttr(parser.aheadToken)              // 设置比较操作符属性
				next.SetLeft(cur)
				parser.match(parser.aheadToken) // 匹配比较符号

				right := parser.additiveExpression()
				next.SetRight(right)
				cur = next
			}
			// cur 作为simple_expression 分析完毕
			res = cur
		default: // 数组变量或单值变量
			cur = NewASTNode(EXPRESSION, VAR, line) // 表达式，左值变量
			cur.SetAttr(id)                         // 设置ID属性
			if parser.aheadToken == L_PARE_M {      // 数组变量
				parser.match(L_PARE_M)
				t := parser.expression()
				cur.SetLeft(t)
				parser.match(R_PARE_M)
			}
			// 判断后续有无'=';
			if parser.aheadToken == ASSIGN { // 赋值语句
				res = NewASTNode(EXPRESSION, ASSIGNMENT, line) // 表达式，赋值语句
				parser.match(ASSIGN)

				t := parser.expression()
				res.SetLeft(cur) // 左子节点为变量
				res.SetRight(t)  //右子节点为表达式
			} else { // var语句
				// cur 作为factor 继续向上分析
				for parser.aheadToken == MUL || parser.aheadToken == DIV {
					next = NewASTNode(EXPRESSION, OPERATION, line) // 表达式，操作符
					next.SetLeft(cur)
					next.SetAttr(parser.aheadToken) // 设置操作符属性
					parser.match(parser.aheadToken) // 匹配操作符

					right := parser.factor()
					next.SetRight(right)
					cur = next
				}
				// cur 作为term 继续向上分析
				for parser.aheadToken == PLUS || parser.aheadToken == MINUS {
					next = NewASTNode(EXPRESSION, OPERATION, line) // 表达式，操作符
					next.SetLeft(cur)
					next.SetAttr(parser.aheadToken) // 设置操作符属性
					parser.match(parser.aheadToken) // 匹配操作符

					right := parser.term()
					next.SetRight(right)
					cur = next
				}
				// cur 作为additive_expression 继续向上分析
				if parser.aheadToken == LE || parser.aheadToken == LT || parser.aheadToken == GT || parser.aheadToken == GE || parser.aheadToken == EQ || parser.aheadToken == NOT_EQ {
					next = NewASTNode(EXPRESSION, COMPARE, line) // 表达式，比较语句
					next.SetAttr(parser.aheadToken)              // 设置比较操作符属性
					next.SetLeft(cur)
					parser.match(parser.aheadToken) // 匹配比较符号

					right := parser.additiveExpression()
					next.SetRight(right)
					cur = next
				}
				// cur 作为simple_expression 分析完毕
				res = cur
			}
		}
	}
	return res
}

// 选择语句
func (parser *Parser) selectionStmt() *ASTNode {
	var res, els *ASTNode
	line := parser.buffer.Lines()

	parser.match(IF)
	parser.match(L_PARE_S)
	exp := parser.expression()
	parser.match(R_PARE_S)
	then := parser.statement()
	if parser.aheadToken == ELSE {
		parser.match(ELSE)
		els = parser.statement()
	}
	res = NewASTNode(STATEMENT, SELECTION_STMT, line) // 语句，选择语句
	res.SetLeft(exp)
	res.SetMid(then)
	res.SetRight(els)
	return res
}

// 循环语句
func (parser *Parser) iterationStmt() *ASTNode {
	var res *ASTNode
	line := parser.buffer.Lines()

	parser.match(WHILE)
	parser.match(L_PARE_S)
	exp := parser.expression()
	parser.match(R_PARE_S)
	then := parser.statement()
	res = NewASTNode(STATEMENT, ITERATION_STMT, line) // 语句，循环语句
	res.SetLeft(exp)
	res.SetMid(then)
	return res
}

// 返回语句
func (parser *Parser) returnStmt() *ASTNode {
	var res, cur *ASTNode
	line := parser.buffer.Lines()

	parser.match(RETURN)
	// 可选的expression部分
	if parser.aheadToken == ID || parser.aheadToken == NUM || parser.aheadToken == L_PARE_S {
		cur = parser.expression()
	}
	parser.match(SEMI)
	res = NewASTNode(STATEMENT, RETURN_STMT, line) // 语句，返回语句
	res.SetLeft(cur)                               // 直接返回则左子节点为空
	return res
}

// 简单表达式，包括加法表达式或关系表达式
func (parser *Parser) simpleExpression() *ASTNode {
	var res *ASTNode
	line := parser.buffer.Lines()

	left := parser.additiveExpression()

	switch parser.aheadToken {
	case LE, LT, GT, GE, EQ, NOT_EQ:
		res = NewASTNode(EXPRESSION, COMPARE, line) // 表达式，比较表达式
		res.SetAttr(parser.aheadToken)              // 设置比较符号
		parser.match(parser.aheadToken)

		res.SetLeft(left)
		right := parser.additiveExpression()
		res.SetRight(right)
	default:
		res = left
	}
	return res
}

// 加法表达式
func (parser *Parser) additiveExpression() *ASTNode {
	var cur, next *ASTNode
	cur = parser.term()
	line := parser.buffer.Lines()

	for parser.aheadToken == PLUS || parser.aheadToken == MINUS {
		next = NewASTNode(EXPRESSION, OPERATION, line) // 操作符表达式
		next.SetAttr(parser.aheadToken)                // 设置操作符
		parser.match(parser.aheadToken)

		t := parser.term()
		next.SetLeft(cur)
		next.SetRight(t)
		cur = next
	}
	return cur
}

// parser.term,乘除法
func (parser *Parser) term() *ASTNode {
	var cur, next *ASTNode
	cur = parser.factor()
	line := parser.buffer.Lines()

	for parser.aheadToken == MUL || parser.aheadToken == DIV {
		next = NewASTNode(EXPRESSION, OPERATION, line) // 操作符表达式
		next.SetAttr(parser.aheadToken)                // 设置操作符
		parser.match(parser.aheadToken)

		t := parser.term()
		next.SetLeft(cur)
		next.SetRight(t)
		cur = next
	}
	return cur
}

// parser.factor，基本元
func (parser *Parser) factor() *ASTNode {
	var res *ASTNode
	var id TokenString
	line := parser.buffer.Lines()

	switch parser.aheadToken {
	case NUM:
		res = NewASTNode(EXPRESSION, CONST, line) // 常量表达式
		val, _ := strconv.ParseInt(string(parser.lexeme), 10, 0)
		res.SetAttr(val) // 设置常量值
		parser.match(NUM)
	case L_PARE_S:
		parser.match(L_PARE_S)
		res = parser.expression()
		parser.match(R_PARE_S)
	case ID: // 左值变量或者函数调用
		id = parser.lexeme
		AddIdentifier(string(parser.lexeme)) // 标识符添加到符号表
		parser.match(ID)
		if parser.aheadToken == L_PARE_S { // 函数调用
			res = NewASTNode(EXPRESSION, CALL, line)
			parser.match(L_PARE_S)
			t := parser.args() // 函数调用实参
			parser.match(R_PARE_S)
			res.SetLeft(t)
		} else { // 左值变量
			res = NewASTNode(EXPRESSION, VAR, line)
			if parser.aheadToken == L_PARE_M { // 数组
				parser.match(L_PARE_M)
				t := parser.expression()
				res.SetLeft(t)
				parser.match(R_PARE_M)
			}
		}
		// 设置ID对应的lexeme
		res.SetAttr(id)
	default:
		parser.syntaxError()
	}
	return res
}

// 函数实参
func (parser *Parser) args() *ASTNode {
	var res, cur, next *ASTNode
	line := parser.buffer.Lines()

	if parser.aheadToken == ID || parser.aheadToken == NUM || parser.aheadToken == L_PARE_S {
		cur = parser.expression()
		res = cur
		for parser.aheadToken == COMMA {
			parser.match(COMMA)
			next = parser.expression()
			cur.SetSibling(next)
			cur = next
		}
	}
	tmp := NewASTNode(ARGS, nil, line) // 实参
	tmp.SetLeft(res)
	return tmp
}

// 利用递归下降法生成抽象语法树
func (parser *Parser) Parse() (*ASTNode, *SymbolTableNode) {
	// 获取第一个token
	for parser.aheadToken, parser.lexeme = parser.scanner.getToken(); parser.aheadToken == COMMENT; parser.aheadToken, parser.lexeme = parser.scanner.getToken() {
	}

	// 初始化当前符号表
	curTable = NewTable()

	// 语法树以声明列表的形式调用
	astNode := parser.declarationList()
	if parser.aheadToken != EOF_TOKEN {
		parser.syntaxError()
	}
	fmt.Println("Parser Done!")
	return astNode, curTable
}

// 匹配期待的token并获取下一个token
func (parser *Parser) match(t Token) {
	if t == parser.aheadToken {
		// 符号表操作
		if t == R_PARE_L {
			MoveUp()
		} else if t == L_PARE_L {
			MoveDown()
		}

		// 获取下一个token,将注释token和错误token过滤
		for parser.aheadToken, parser.lexeme = parser.scanner.getToken(); parser.aheadToken == COMMENT || parser.aheadToken == ERROR; parser.aheadToken, parser.lexeme = parser.scanner.getToken() {
		}
	} else {
		parser.syntaxError()
	}
}

// 语法错误时打印错误消息
func (parser *Parser) syntaxError() {
	fmt.Printf("%s: [%d]. Token [%d]\n", "Syntax Error in Line", parser.buffer.Lines(), parser.aheadToken)
	// 获取下一个token,将注释token和错误token过滤
	for parser.aheadToken, parser.lexeme = parser.scanner.getToken(); parser.aheadToken == COMMENT || parser.aheadToken == ERROR; parser.aheadToken, parser.lexeme = parser.scanner.getToken() {
	}
}

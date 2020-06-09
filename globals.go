// Copyright 2020. All rights reserved.
// Author: Zhifei Liu, 2020/6
// Filename: globals.go
// Package: scan
// Description: 本文件定义了词法分析器和语法分析器使用的常量

package scan

import "os"

var (
	FileOut *os.File // 输出的文件指针

	BufferConst  *Buffer  // 输入缓冲区
	ScannerConst *Scanner // 词法扫描器
	ParserConst  *Parser  // 语法分析器
)

// 符号表
// 扫描阶段即可开始设置其中的一些值
// 符号表移动策略，遇到{移到下一层，遇到}移到上一层
// match函数处进行符号表的基本操作
var curTable, nextTable, prevTable *SymbolTableNode
var hasSiblings bool = false // 判断是否有同层次的并列的scope

// 词法分析needed
type Token int          // 扫描的token类型
type TokenString []byte // 扫描的词素类型
type StateType int      // DFA状态类型

// 语法分析needed
type NodeKind int // 节点类型,语句、表达式、形参、实参、类型等。
type StmtKind int // 语句类型
type ExpKind int  // 表达式类型
type ExpType int  // 表达式值类型,int、void、bool等
type VarType int  // 形参、实参、变量声明、函数返回类型等

// 节点类型常量
const (
	STATEMENT  NodeKind = iota // 语句类型
	EXPRESSION                 // 表达式类型
	PARAMS                     // 形式参数列表
	PARAM                      // 单个参数
	ARGS                       // 实参
	TYPE                       // 类型
)

// 语句子类型
const (
	VAR_DECLARATION  StmtKind = iota // 变量声明语句
	FUNC_DECLARATION                 // 函数声明语句
	COMPOUND                         // 复合语句
	SELECTION_STMT                   // 选择语句
	ITERATION_STMT                   // 循环语句
	RETURN_STMT                      // 返回语句
)

// 表达式子类型
const (
	VAR        ExpKind = iota // 左值变量
	ASSIGNMENT                // 赋值语句
	CALL                      // 函数调用
	COMPARE                   // 比较语句
	CONST                     // 常数节点
	OPERATION                 // 操作符节点
)

// 表达式值类型常量
const (
	EXP_INT ExpType = iota
	EXP_BOOL
	EXP_VOID
)

// 变量类型
const (
	VAR_TYPE_INT        VarType = iota // int
	VAR_TYPE_INT_VECTOR                // int数组
	VAR_TYPE_VOID                      // void
)

// DFA状态
const (
	// 开始结束状态
	START StateType = iota
	DONE
	//
	INID
	INNUM
	INLT
	INGT
	INEQ
	NOT
	// 注释
	INCOM_B
	INCOM_C
	INCOM_D
)

// 指示字符EOF
const (
	EOF_CHAR byte = 0
)

// Token Type
const (
	// Key words
	IF Token = iota
	ELSE
	INT
	RETURN
	VOID
	WHILE

	// OP
	PLUS
	MINUS
	MUL
	DIV
	LT
	LE
	GT
	GE
	EQ
	NOT_EQ
	ASSIGN

	// other
	SEMI
	COMMA
	L_PARE_S
	L_PARE_M
	L_PARE_L
	R_PARE_S
	R_PARE_M
	R_PARE_L

	// ID, NUM
	ID
	NUM

	// ERROR and EOF
	COMMENT
	ERROR
	EOF_TOKEN
)

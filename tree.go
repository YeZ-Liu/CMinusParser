// Copyright 2020. All rights reserved.
// Author: Zhifei Liu, 2020/6
// Filename: tree.go
// Package: scan
// Description: 本文件定义了抽象语法树节点类以及成员函数和初始化工厂函数

package scan

// 抽象语法树节点
type ASTNode struct {
	line             int         // 节点所处行号
	nodeK            NodeKind    //节点类型(四种)
	nodeT            interface{} // 具体节点类型(如变量、选择语句等)
	attribute        interface{} // 不同节点的不同属性,如词素、值、操作符、数组大小等
	varT             VarType     // 变量类型、形参类型、函数返回类型等
	expT             ExpType     // 表达式结果类型,用于后续类型检查
	left, right, mid *ASTNode    // 子节点
	sibling          *ASTNode    // 兄弟节点
}

// 设置左右中子树
func (node *ASTNode) SetLeft(son *ASTNode) {
	node.left = son
}

func (node *ASTNode) SetRight(son *ASTNode) {
	node.right = son
}
func (node *ASTNode) SetMid(son *ASTNode) {
	node.mid = son
}

// 设置兄弟节点
func (node *ASTNode) SetSibling(s *ASTNode) {
	node.sibling = s
}

// 设置表达式结果类型
func (node *ASTNode) SetExpType(t ExpType) {
	node.expT = t
}

// 设置变量类型,第一个参数是类型
func (node *ASTNode) SetType(t Token) {
	switch t {
	case INT:
		node.varT = VAR_TYPE_INT
	case VOID:
		node.varT = VAR_TYPE_VOID
	}
}

// 设置为数组
func (node *ASTNode) SetVec() {
	node.varT = VAR_TYPE_INT_VECTOR
}

// 设置其他属性
func (node *ASTNode) SetAttr(attr interface{}) {
	node.attribute = attr
}

// 结构体工厂函数，返回一个节点指针
func NewASTNode(k NodeKind, t interface{}, l int) *ASTNode {
	var newNode ASTNode
	newNode = ASTNode{nodeK: k, nodeT: t, line: l}
	return &newNode
}

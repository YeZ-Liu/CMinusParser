// Copyright 2020. All rights reserved.
// Author: Zhifei Liu, 2020/6
// Filename: symbolTable.go
// Package: scan
// Description: 本文件定义了符号表类和相应的成员函数和工厂函数
// 				这个文件并没有实现完整，如果需要进行语法制导翻译则需要好好完善
// 				暂时只将扫描到的词素保存到相应域的符号表中

package scan

import (
	"errors"
)

// 链接符号表节点，每一个域都有一个符号表
// 数据部分由map构成，键为标识符的lexeme，值为接口类型(方便自定义各种属性)
// 同时包含多个指针指向下一层或上一层或同层的符号表
type SymbolTableNode struct {
	data                          map[string]interface{}
	prev, next, rightSib, leftSib *SymbolTableNode
}

// 新建一个符号表节点，初始化数据部分最多10个标识符
func NewTable() *SymbolTableNode {
	var res SymbolTableNode
	res.data = make(map[string]interface{}, 10)
	return &res
}

// 设置下一层指针
func (node *SymbolTableNode) SetNext(n *SymbolTableNode) {
	node.next = n
}

// 设置上一层指针
func (node *SymbolTableNode) SetPrev(n *SymbolTableNode) {
	node.prev = n
}

// 设置兄弟节点指针
func (node *SymbolTableNode) SetRightSibling(n *SymbolTableNode) {
	node.rightSib = n
}

// 设置兄弟节点指针
func (node *SymbolTableNode) SetLeftSibling(n *SymbolTableNode) {
	node.leftSib = n
}

// 返回下一层的符号表指针
func (node *SymbolTableNode) Next() *SymbolTableNode {
	if node.next != nil {
		right := node.next
		for ; right.rightSib != nil; right = right.rightSib {
		}
		return right
	}
	return nil
}

// 返回上一层的符号表指针
func (node *SymbolTableNode) Prev() *SymbolTableNode {
	if node.prev != nil {
		return node.prev
	}
	return nil
}

// 返回兄弟节点符号表指针
func (node *SymbolTableNode) LeftSibling() *SymbolTableNode {
	if node.leftSib != nil {
		return node.leftSib
	}
	return nil
}

// 返回兄弟节点符号表指针
func (node *SymbolTableNode) RightSibling() *SymbolTableNode {
	if node.rightSib != nil {
		return node.rightSib
	}
	return nil
}

// 添加一条新的标识符属性到符号表
func (node *SymbolTableNode) Put(id string, c interface{}) error {
	if _, ok := node.data[id]; ok {
		return errors.New("the ID has existed")
	}

	node.data[id] = c
	return nil
}

// 获取相应标识符的属性
func (node *SymbolTableNode) Get(key string) interface{} {
	if val, ok := node.data[key]; ok {
		return val
	}
	return nil
}

// 将符号表向下移动
func MoveDown() {
	nextTable = curTable.Next()          // 下一层的最右节点
	if hasSiblings && nextTable != nil { // 下一节点必然非空
		// 创建新的同层节点
		newNode := NewTable()
		nextTable.SetRightSibling(newNode)
		newNode.SetPrev(curTable)
		nextTable = newNode
	} else { // 下一节点必然为空
		nextTable = NewTable()
		nextTable.SetPrev(curTable)
		curTable.SetNext(nextTable)
	}
	curTable = nextTable // 重置当前节点
	hasSiblings = false
}

// 将符号表向上移动
func MoveUp() {
	hasSiblings = true
	curTable = curTable.Prev()
	if curTable == nil {
		ParserConst.syntaxError()
	}
}

// 向符号表添加标识符
func AddIdentifier(lexeme string) {
	err := curTable.Put(lexeme, NewContent(BufferConst.Lines()))
	if err != nil {
		//fmt.Println(err.Error(), ":", lexeme)
	}
}

// 标识符属性域
// 标识符分为函数名、变量名、形参名。。。
// 函数标识符属性：参数列表，返回值类型，
// 变量名属性：类型，是否数组，数组大小
// 形参属性：类型，是否数组
type content struct {
	line   int      // 变量或函数第一次出现的置行号
	kind_  int      //变量或是函数或是其他，1变量，2函数
	type_  Token    // 变量类型，函数返回类型
	size_  int64    // 是否是数组变量(>0是)
	params []string // 函数类型标识符的参数列表对应标识符
}

// 返回属性域
// 入口参数：行号
func NewContent(l int) *content {
	return &content{line: l}
}

// 添加其他属性
// 种类，数组大小，变量类型或返回值类型
func (c *content) AddAttr(k int, s int64, t Token) {
	c.type_ = t
	c.size_ = s
	c.kind_ = k
}

// 为函数类型标识符添加参数标识符
func (c *content) AddParam(param string) {
	c.params = append(c.params, param)
}

// Copyright 2020. All rights reserved.
// Author: Zhifei Liu, 2020/6
// Filename: buffer.go
// Package: scan
// Description: 本文件定义了词法分析器使用的输入缓冲区类及相应的成员函数和工厂函数

package scan

import (
	"bufio"
	"io"
	"log"
	"os"
)

// 输入缓冲区类
// 使用缓冲区读取文件内容并记录行数
type Buffer struct {
	file   *os.File      // 读取的文件指针
	reader *bufio.Reader // 缓冲区
	line   int           // 当前读取字符行号
}

// 获取指定文件名缓冲区
func NewBuffer(file string) *Buffer {
	var b Buffer
	var err error
	b.file, err = os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	b.reader = bufio.NewReader(b.file)
	b.line = 1
	return &b
}

// 获取缓冲区输入当前行数
func (b *Buffer) Lines() int {
	return b.line
}

// 获取下一个字符
func (b *Buffer) Next() (res byte) {
	var err error
	res, err = b.reader.ReadByte()
	if err == io.EOF {
		res = EOF_CHAR
		b.file.Close() // 读取完后关闭文件
	}
	if res == '\n' {
		b.line++
	}
	return
}

// 回退字符指针
// 回退一次后再读取再回退，以此确认换行符
func (b *Buffer) UnNext() {
	var err error
	err = b.reader.UnreadByte()
	if err != nil {
		log.Println("Can't back up!")
		return
	}
	cur, err := b.reader.ReadByte()
	if cur == '\n' {
		b.line--
	}
	err = b.reader.UnreadByte()
}

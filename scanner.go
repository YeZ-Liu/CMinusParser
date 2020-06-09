// Copyright 2020. All rights reserved.
// Author: Zhifei Liu, 2020/6
// Filename: scanner.go
// Package: scan
// Description: 本文件定义了词法分析器类以及成员函数和初始化工厂函数

package scan

import (
	"unicode"
)

// 词法分析器类
type Scanner struct {
	KeyTable map[string]Token  // 关键字表
	buffer *Buffer             // 输入缓冲区


}

// 初始化关键字表
func (scanner *Scanner) initKeyTable()  {
	scanner.KeyTable = make(map[string]Token, 6)
	scanner.KeyTable["if"] = IF
	scanner.KeyTable["else"] = ELSE
	scanner.KeyTable["int"] = INT
	scanner.KeyTable["void"] = VOID
	scanner.KeyTable["return"] = RETURN
	scanner.KeyTable["while"] = WHILE
}

/**
	判断ID token 是够为关键字
	如果是,则返回对应的关键字类型
	否则,返回ID
 */
func (scanner *Scanner) idToken(s string) Token {
	if token, ok := scanner.KeyTable[s]; ok {
		return token
	}
	return ID
}

// 词法分析器类工厂函数
// 参数为输入文件缓冲区
func NewScanner(buf *Buffer) *Scanner  {
	var scanner Scanner
	scanner.buffer = buf   // 初始化输入缓冲区
	scanner.initKeyTable()  // 初始化关键字表
	return &scanner
}


// 从输入缓冲中扫描token,返回Token类型和Token的词素
func (scanner *Scanner) getToken() (Token, TokenString) {
	state := START         // DFA开始状态
	var lexeme TokenString // 扫描到的词素
	var token Token        // 扫描到的token
	var char byte          // 输入的下一个字符
	var save bool          // 是否保存当前扫描的字符(可能回退，空白符等)

	for state != DONE {
		save = true

		// 读取下一个字符
		char = scanner.buffer.Next()
		if char == EOF_CHAR { // 文件结尾，返回EOF Token
			token = EOF_TOKEN
			break
		}

		switch state {
		case START:
			switch {
			case unicode.IsSpace(rune(char)): // 空白符
				state = START
				save = false
			case unicode.IsLetter(rune(char)): // 字符
				state = INID
			case unicode.IsDigit(rune(char)): // 数字
				state = INNUM
			case char == '<':
				state = INLT
			case char == '>':
				state = INGT
			case char == '=':
				state = INEQ
			case char == '!':
				state = NOT
			case char == '/':
				state = INCOM_B
			case char == ';':
				token = SEMI
				state = DONE
			case char == ',':
				token = COMMA
				state = DONE
			case char == '(':
				token = L_PARE_S
				state = DONE
			case char == ')':
				token = R_PARE_S
				state = DONE
			case char == '[':
				token = L_PARE_M
				state = DONE
			case char == ']':
				token = R_PARE_M
				state = DONE
			case char == '{':
				token = L_PARE_L
				state = DONE
			case char == '}':
				token = R_PARE_L
				state = DONE
			case char == '+':
				token = PLUS
				state = DONE
			case char == '-':
				token = MINUS
				state = DONE
			case char == '*':
				token = MUL
				state = DONE
			default: // 默认其他输入错误
				state = DONE
				token = ERROR
			}
		case INID:
			if !unicode.IsLetter(rune(char)) { // 读取到ID，回退当前字符
				scanner.buffer.UnNext()
				token = ID
				state = DONE
				save = false
			} else {
				state = INID
			}
		case INNUM:
			if !unicode.IsNumber(rune(char)) { // 读取到NUMBER，回退当前字符
				scanner.buffer.UnNext()
				token = NUM
				state = DONE
				save = false
			} else {
				state = INNUM
			}
		case INLT:
			if char == '=' {
				token = LE
				state = DONE
			} else {
				scanner.buffer.UnNext()
				token = LT
				state = DONE
				save = false
			}
		case INGT:
			if char == '=' {
				token = GE
				state = DONE
			} else {
				scanner.buffer.UnNext()
				token = GT
				state = DONE
				save = false
			}
		case INEQ:
			if char == '=' {
				token = EQ
				state = DONE
			} else {
				scanner.buffer.UnNext()
				token = ASSIGN
				state = DONE
				save = false
			}
		case NOT:
			if char == '=' {
				token = NOT_EQ
				state = DONE
			} else {
				scanner.buffer.UnNext()
				token = ERROR
				state = DONE
				save = false
			}
		case INCOM_B:
			if char == '*' { // 注释，默认保存注释lexeme
				state = INCOM_C
			} else { // 除号，回退当前字符
				scanner.buffer.UnNext()
				token = DIV
				state = DONE
				save = false
			}
		case INCOM_C:
			if char == '*' {
				state = INCOM_D
			} else {
				state = INCOM_C
			}
		case INCOM_D:
			if char == '*' {
				state = INCOM_D
			} else if char == '/' { // 完成注释扫描
				token = COMMENT
				state = DONE
			} else {
				state = INCOM_C
			}
			// 只有上述DFA状态，因此没有默认状态处理
		}

		// 是否需要将当前字符保存到lexeme
		if save {
			lexeme = append(lexeme, char)
		}
	}

	// 在关键字表里查找当前扫描ID是否是关键字
	if token == ID {
		token = scanner.idToken(string(lexeme))
	}

	// 调用辅助函数打印结果
	//HelpPrint(token, lexeme, scanner.buffer.Lines())
	HelpPrintFile(token, lexeme, scanner.buffer.Lines(), FileOut)

	return token, lexeme
}

// 只进行词法扫描
func (scanner *Scanner) ScanAll()  {
	// 获取token和词素并打印
	for token, tokenString := scanner.getToken(); token != EOF_TOKEN; token, tokenString = scanner.getToken() {
		HelpPrintFile(token, tokenString, scanner.buffer.Lines(), FileOut)
	}
}
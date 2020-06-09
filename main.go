package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"scan"
)

var (
	f    string
	v, V bool
	h    bool
	s, p bool
	c  bool
)

func init() {
	flag.BoolVar(&h, "h", false, "帮助信息")
	flag.BoolVar(&v, "v", false, "版本信息")
	flag.BoolVar(&V, "V", false, "版本信息")
	flag.BoolVar(&s, "s", false, "词法分析")
	flag.BoolVar(&p, "p", false, "语法分析")
	flag.BoolVar(&c, "c", false, "标准输出")

	flag.StringVar(&f, "f", "", "`filename`")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `CMinusParser version: CMinusParser/1.0.1
Usage: CMinusParser [-hvV] -sp -f filename

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	if v || V {
		fmt.Fprintf(os.Stderr, `CMinusParser version: CMinusParser/1.0.1`)
		return
	}

	if len(f) == 0 {
		fmt.Println("请输入文件完整路径名!")
		return
	}

	filename := f

	// 获取文件名字
	stat, err := os.Stat(filename)
	if err != nil {
		fmt.Println("输入文件有误!")
		return
	}
	name := stat.Name()

	// 获取输入文件路径
	dir := strings.TrimRight(filename, name)

	//fmt.Println(dir)  // C:/Users/lzff1/Desktop/


	newName := dir + "CMinusParserOut.txt"

	// 先将文件删除
	err = os.Remove(newName)
	if err != nil {
		// 删除失败不需要提示
	}

	// 输出到标准输出
	if c {
		scan.FileOut = os.Stdout
	} else {
		// 新建一个文件作为输入
		scan.FileOut, err = os.OpenFile(newName, os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			fmt.Println("输出文件创建失败!")
			return
		}
	}


	// 初始化缓冲区
	scan.BufferConst = scan.NewBuffer(filename)

	// 语法分析
	if p {
		scan.ParserConst = scan.NewParser(scan.BufferConst)
		astRoot, tableRoot := scan.ParserConst.Parse()
		scan.HelpPrintTree(astRoot, 0, '-', scan.FileOut)
		scan.HelpPrintTable(tableRoot, 0, '-',  scan.FileOut)
	} else if s {
		// 只进行词法分析
		scan.ScannerConst = scan.NewScanner(scan.BufferConst)
		scan.ScannerConst.ScanAll()
	} else {
		flag.Usage()
	}
}

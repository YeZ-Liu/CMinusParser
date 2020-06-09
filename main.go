package scan

import (
	"flag"
	"fmt"
	"os"
)

var (
	f    string
	v, V bool
	h    bool
	s, p bool
)

func init() {
	flag.BoolVar(&h, "h", false, "帮助信息")
	flag.BoolVar(&v, "v", false, "版本信息")
	flag.BoolVar(&V, "V", false, "版本信息")
	flag.BoolVar(&s, "s", false, "词法分析")
	flag.BoolVar(&p, "p", false, "语法分析")
	flag.StringVar(&f, "f", "", "`文件名`")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `CMinusParser version: CMinusParser/1.0.1
Usage: CMinusParser [-hvVsp] [-f filename]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

}

package main

import (
	"apidoc/option"
	"github.com/jessevdk/go-flags"
	_ "github.com/jessevdk/go-flags"
)

func main() {
	p := flags.NewParser(&option.Opt, flags.Default)
	p.ShortDescription = "文档生成"
	p.Parse()
}

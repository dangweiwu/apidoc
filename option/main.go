package option

import (
	"fmt"
	"github.com/dangweiwu/apidoc/internal/filex"
	"github.com/dangweiwu/apidoc/internal/parse"
	"os"
	"path/filepath"
)

var Opt Option

type Option struct {
	Root Root `command:"run" description:"生成文档"`
}

type Root struct {
	Root string `short:"r" description:"根路径"`
	Name string `short:"o" description:"输出文件名"`
}

func (this *Root) Usage() string {
	return `
api workdown文档生成
默认根目录当前文件夹
`
}
func (this *Root) Execute(args []string) error {
	var (
		err  error
		root string
	)
	root = Opt.Root.Root
	if len(root) == 0 {
		if root, err = os.Getwd(); err != nil {
			panic(err)
		}
	}
	fmt.Printf("[根目录] %s\n", root)

	fobj := filex.NewFilex(root)
	//fobj := &filex.Filex{Root: root}
	err = fobj.GetModule()
	if err != nil {
		fmt.Printf("[ERR]: %s\n", err)
		return nil
	}

	fmt.Printf("[MODULE]: %s\n", fobj.Module)

	fobj.Walk()
	if len(this.Name) == 0 {
		os.MkdirAll(filepath.Join(this.Root, "doc"), 0777)
		this.Name = filepath.Join(this.Root, "doc", "api.md")
	}

	if err := parse.NewMdOut(fobj.Parse.Doc, this.Name).Out(); err == nil {
		fmt.Printf("[ok]: 生成文件:%s \n", this.Name)
	} else {
		fmt.Printf("[err]: 生成失败 :%v\n", err)
	}

	return nil
}

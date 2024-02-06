package option

import (
	"apidoc/internal/filex"
	"fmt"
	"os"
)

var Opt Option

type Option struct {
	Root Root `command:"run" description:"生成文档"`
}

type Root struct {
	Root string `short:"r" description:"根路径"`
}

func (this *Root) Usage() string {
	return `
	api workdown文档生成`
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

	fobj := &filex.Filex{Root: root}
	err = fobj.GetModule()
	if err != nil {
		fmt.Printf("[ERR]: %s\n", err)
		return nil
	}

	fmt.Printf("[MODULE]: %s\n", fobj.Module)

	fobj.Walk()

	return nil
}

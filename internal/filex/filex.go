package filex

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Filex struct {
	Root   string
	Module string
}

func (this *Filex) GetModule() error {
	if len(this.Root) == 0 {
		return errors.New("need path root")
	}

	modfile := path.Join(this.Root, "go.mod")

	f, err := os.Open(modfile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("read go.mod error:%w", err)
	}
	buf := bufio.NewReader(f)
	moduleLinebts, _, err := buf.ReadLine()
	if err != nil {
		return fmt.Errorf("read go.mod line error:%w", err)
	}

	this.Module = strings.TrimSpace(strings.Replace(string(moduleLinebts), "module", "", -1))
	return nil

}

func (this *Filex) visitFile(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if fi.IsDir() {
		fmt.Printf("[dir] %s %s\n", fp, fi.Name())
		return nil
	}
	fmt.Printf("[FILE]: %s %s \n", fp, fi.Name())
	return nil
}
func (this *Filex) Walk() {
	filepath.Walk(this.Root, this.visitFile)
}

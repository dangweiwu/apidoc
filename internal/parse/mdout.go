package parse

import (
	_ "embed"
	"os"
	"text/template"
)

var (
	//go:embed mdtpl.tpl
	md_tpl string
)

// markdown 格式输出
type MdOut struct {
	DocData *DocData
	name    string
}

func NewMdOut(d *DocData, name string) *MdOut {
	return &MdOut{d, name}
}

func (this *MdOut) Out() error {
	tpl, err := template.New("doc").Parse(md_tpl)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(this.name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	defer f.Close()
	if err != nil {
		return err
	}

	err = tpl.Execute(f, this.DocData)
	return err
}

package generate

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed assets/README.md.tmpl
var readmeTmpl string

type readmeModel struct {
	Name      string
}

func (g generator) GenerateReadme() error {
	fmt.Println("Generating README.md...")
	m := readmeModel{
		Name: g.modulePath.Package,
	}

	tmpl, err := template.New("readme").Parse(readmeTmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(g.outPath, "/README.md")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, m)
}

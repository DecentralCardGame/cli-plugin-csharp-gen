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
//go:embed assets/gitignore
var gitIgnore []byte

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

	path := filepath.Join(g.outPath, "README.md")
	readme, err := os.Create(path)
	if err != nil {
		return err
	}
	defer readme.Close()

	err = os.WriteFile(filepath.Join(g.outPath, ".gitignore"), gitIgnore, os.ModePerm)
	if err != nil {
		return err
	}

	return tmpl.Execute(readme, m)
}

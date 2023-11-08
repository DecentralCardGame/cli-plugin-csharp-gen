package generate

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed assets/csproj.tmpl
var csprojTmpl string

type model struct {
	Name      string
	ShortName string
	URL       string
}

func (g generator) GenerateCsproj() error {
	fmt.Println("Generating csproj...")
	m := model{
		URL:       "https://" + g.modulePath.RawPath,
		Name:      strings.Join(strings.Split(g.modulePath.RawPath, "/")[1:], "."),
		ShortName: g.modulePath.Package,
	}

	tmpl, err := template.New("csproj").Parse(csprojTmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(g.outPath, "/"+g.modulePath.Package+".csproj")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, m)
}

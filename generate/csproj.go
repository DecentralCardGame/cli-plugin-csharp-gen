package generate

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed assets/csproj.tmpl
var csprojTmpl string

type csprojModel struct {
	Name      string
	ShortName string
	URL       string
}

func (g generator) GenerateCsproj() error {
	fmt.Println("Generating csproj...")
	m := csprojModel{
		URL:       "https://" + g.modulePath.RawPath,
		Name:      g.csModulePath,
		ShortName: g.modulePath.Package,
	}

	tmpl, err := template.New("csproj").Parse(csprojTmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(g.outPath, g.csModulePath+".csproj")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, m)
}

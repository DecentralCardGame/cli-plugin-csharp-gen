package generate

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed assets/Client.cs.tmpl
var clientTmpl string

type ClientModel struct {
	Path    string
	Txs     []ServiceModel
	Queries []ServiceModel
}

type ServiceModel struct {
	Path string
	Name string
	Type string
}

func (g generator) GenerateClient(ctx context.Context) error {
	fmt.Println("Generating client...")
	cache := protoanalysis.NewCache()
	pkgs, err := protoanalysis.Parse(ctx, cache, g.protoPath)
	if err != nil {
		return err
	}

	model := ClientModel{
		Path: g.csModulePath,
	}

	for _, pkg := range pkgs {
		for _, service := range pkg.Services {
			s := ServiceModel{
				Type: service.Name,
				Path: strings.Title(pkg.Name),
				Name: getSimpleModuleNameFromPath(pkg.Name),
			}
			switch service.Name {
			case "Query":
				model.Queries = append(model.Queries, s)
			case "Msg":
				model.Txs = append(model.Txs, s)
			}
		}
	}
	
	tmpl, err := template.New("client").Parse(clientTmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(g.outPath, "Client.cs")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, model)
}

func getSimpleModuleNameFromPath(path string) string {
	s := strings.Split(path, ".")
	
	if strings.Contains(s[len(s) - 1], "v1") {
		return strings.Title(s[len(s) - 2]) + strings.Title(s[len(s) - 1]) 
	} else {
		return strings.Title(s[len(s) - 1])
	}
}

package generate

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/ignite/cli/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
	"path/filepath"
	"strings"
)

//go:embed assets/buf.gen.grpc.yaml
var bufGenGrpcYaml string

type QueryFile struct {
	name string
	module string
}

func (g generator) GenerateQueries(ctx context.Context) error {
	fmt.Println("Generating queries...")
	buf, err := cosmosbuf.New()
	if err != nil {
		return err
	}
	defer buf.Cleanup()

	cache := protoanalysis.NewCache()
	var queryFiles []QueryFile
	pkgs, err := protoanalysis.Parse(ctx, cache, g.protoPath)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			if strings.Contains(file.Path, "/query.proto") || strings.Contains(file.Path, "/service.proto") {
				queryFiles = append(queryFiles, QueryFile{file.Path, strings.Replace(strings.Title(pkg.Name), ".", "/", -1)})
			}
		}
	}
	
	for _, file := range queryFiles {
		err := buf.Generate(ctx, file.name, filepath.Join(g.outPath, file.module), bufGenGrpcYaml)
		if err != nil {
			return err
		}
	}
	
	return nil
}

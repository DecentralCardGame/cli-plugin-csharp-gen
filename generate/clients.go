package generate

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/ignite/cli/ignite/pkg/cosmosbuf"
)

//go:embed assets/buf.gen.yaml
var bufGenYaml string

func (g generator) GenerateClients(ctx context.Context) error {
	fmt.Println("Generating clients...")
	buf, err := cosmosbuf.New()
	if err != nil {
		return err
	}
	
	defer buf.Cleanup()
	return buf.Generate(ctx, g.protoPath, g.outPath, bufGenYaml)
}
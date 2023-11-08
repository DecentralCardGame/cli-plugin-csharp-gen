package main

import (
	"context"
	_ "embed"
	"encoding/gob"
	"fmt"

	"cli-plugin-csharp-gen/generate"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/ignite/services/plugin"
)

func init() {
	gob.Register(plugin.Manifest{})
	gob.Register(plugin.ExecutedCommand{})
	gob.Register(plugin.ExecutedHook{})
}

type p struct{}

func (p) Manifest() (plugin.Manifest, error) {
	return plugin.Manifest{
		Name: "csharp",
		// Add commands here
		Commands: []plugin.Command{
			// Example of a command
			{
				Use:   "csharp-client",
				Short: "Generates csharp client",
				Long:  "Generates csharp client",
				Flags: []plugin.Flag{
					{Name: "out", Type: plugin.FlagTypeString, Usage: "csharp output directory"},
				},
				PlaceCommandUnder: "generate",
			},
		},
		// Add hooks here
		Hooks: []plugin.Hook{},
	}, nil
}

func (p) Execute(cmd plugin.ExecutedCommand) error {
	ctx := context.Background()

	gen, err := generate.New(cmd)
	if err != nil {
		return err
	}

	err = gen.GenerateClients(ctx)
	if err != nil {
		return err
	}

	err = gen.GenerateCsproj()
	if err != nil {
		return err
	}

	return nil
}

func (p) ExecuteHookPre(hook plugin.ExecutedHook) error {
	fmt.Printf("Executing hook pre %q\n", hook.Name)
	return nil
}

func (p) ExecuteHookPost(hook plugin.ExecutedHook) error {
	fmt.Printf("Executing hook post %q\n", hook.Name)
	return nil
}

func (p) ExecuteHookCleanUp(hook plugin.ExecutedHook) error {
	fmt.Printf("Executing hook cleanup %q\n", hook.Name)
	return nil
}

func main() {
	pluginMap := map[string]hplugin.Plugin{
		"cli-plugin-csharp-gen": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}

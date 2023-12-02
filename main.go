package main

import (
	"context"
	_ "embed"
	"encoding/gob"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"cli-plugin-csharp-gen/generate"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/ignite/services/plugin"
)

type Component string

const (
	Component_Clients Component = "clients"
	Component_Csproj  Component = "csproj"
	Component_Readme  Component = "readme"
	Component_Queries Component = "queries"
	Component_Client  Component = "client"
)

func Component_values() []Component {
	return []Component{Component_Clients, Component_Csproj, Component_Readme, Component_Queries, Component_Client}
}

func Component_stringValues() (stringValues []string) {
	for _, val := range Component_values() {
		stringValues = append(stringValues, string(val))
	}
	return
}

func getBuildComponents(cmd *plugin.ExecutedCommand) (components []Component, err error) {
	flags, err := cmd.NewFlags()
	if err != nil {
		return nil, err
	}

	rawComponents, _ := flags.GetStringSlice("components")
	if len(rawComponents) == 0 {
		return Component_values(), nil
	}
	for _, comp := range rawComponents {
		if !slices.Contains(Component_stringValues(), comp) {
			err = fmt.Errorf("buildcomponent '%s' does not exist; options are: [%s]", comp, strings.Join(Component_stringValues(), ", "))
			return
		}
		components = append(components, Component(comp))
	}
	return
}

const cosmosCsharpPluginVersion = "0.1.0"

func init() {
	gob.Register(plugin.Manifest{})
	gob.Register(plugin.ExecutedCommand{})
	gob.Register(plugin.ExecutedHook{})
}

type p struct{}

func (p) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name: "csharp",
		// Add commands here
		Commands: []*plugin.Command{
			// Example of a command
			{
				Use:   "csharp-client",
				Short: "Generates csharp client",
				Long:  "Generates csharp client",
				Flags: []*plugin.Flag{
					{Name: "out", Type: plugin.FlagTypeString, Usage: "csharp output directory"},
					{
						Name: "components",
						Type: plugin.FlagTypeStringSlice,
						Usage: fmt.Sprintf(
							"components to be generated; options: [%s]",
							strings.Join(Component_stringValues(), ", "),
						),
					},
				},
				PlaceCommandUnder: "generate",
			},
		},
		// Add hooks here
		Hooks: []*plugin.Hook{},
	}, nil
}

func (p) Execute(ctx context.Context, cmd *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	buildComponents, err := getBuildComponents(cmd)
	if err != nil {
		return fmt.Errorf("error while getting build components: %s", err.Error())
	}

	err = installPlugin()
	if err != nil {
		return err
	}

	gen, err := generate.New(cmd)
	if err != nil {
		return err
	}

	var componentRegister = map[Component]func() error{
		Component_Csproj: gen.GenerateCsproj,
		Component_Readme: gen.GenerateReadme,
		Component_Clients: func() error {
			return gen.GenerateClients(ctx)
		},
		Component_Queries: func() error {
			return gen.GenerateQueries(ctx)
		},
		Component_Client: func() error {
			return gen.GenerateClient(ctx)
		},
	}

	for _, comp := range buildComponents {
		err = componentRegister[comp]()
		if err != nil {
			return err
		}
	}

	return gen.Build()
}

func installPlugin() error {
	fmt.Println("Installing protoc plugin...")
	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("go install \"github.com/DecentralCardgame/protoc-gen-cosmosCsharp@v%s\"", cosmosCsharpPluginVersion),
	)
	return cmd.Run()
}

func (p) ExecuteHookPre(_ context.Context, hook *plugin.ExecutedHook, api plugin.ClientAPI) error {
	return nil
}

func (p) ExecuteHookPost(_ context.Context, hook *plugin.ExecutedHook, api plugin.ClientAPI) error {
	return nil
}

func (p) ExecuteHookCleanUp(_ context.Context, hook *plugin.ExecutedHook, api plugin.ClientAPI) error {
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"cli-plugin-csharp-gen": plugin.NewGRPC(&p{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}

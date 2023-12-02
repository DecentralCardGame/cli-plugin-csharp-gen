package generate

import (
	chainConfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/services/chain"
	"github.com/ignite/cli/ignite/services/plugin"
	"path/filepath"
	"strings"
)

type generator struct {
	modulePath   gomodulepath.Path
	config       *chainConfig.Config
	appPath      string
	protoPath    string
	outPath      string
	csModulePath string
}

func New(cmd *plugin.ExecutedCommand) (*generator, error) {
	flags, err := cmd.NewFlags()
	if err != nil {
		return nil, err
	}

	c, err := getChain(cmd)
	if err != nil {
		return nil, err
	}
	config, err := c.Config()
	if err != nil {
		return nil, err
	}

	p, appPath, err := getPath(cmd)
	if err != nil {
		return nil, err
	}

	outFlag, _ := flags.GetString("out")
	if outFlag == "" {
		outFlag = "./cs"
	}

	csModulePath := getModulePath(p.RawPath)

	gen := generator{
		config:       config,
		modulePath:   p,
		appPath:      appPath,
		protoPath:    filepath.Join(appPath, config.Config.Build.Proto.Path),
		outPath:      filepath.Join(appPath, outFlag),
		csModulePath: csModulePath,
	}

	return &gen, nil
}

func getChain(cmd *plugin.ExecutedCommand, chainOption ...chain.Option) (*chain.Chain, error) {
	flags, err := cmd.NewFlags()
	if err != nil {
		return nil, err
	}

	var (
		home, _ = flags.GetString("home")
		path, _ = flags.GetString("path")
	)
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return chain.New(absPath, chainOption...)
}

func getPath(cmd *plugin.ExecutedCommand) (gomodulepath.Path, string, error) {
	flags, err := cmd.NewFlags()
	if err != nil {
		return gomodulepath.Path{}, "", err
	}
	path, _ := flags.GetString("path")
	absPath, err := filepath.Abs(path)
	if err != nil {
		return gomodulepath.Path{}, "", err
	}

	return gomodulepath.Find(absPath)
}

func getModulePath(rawPath string) string {
	return strings.Join(strings.Split(rawPath, "/")[1:], ".")
}

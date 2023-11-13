package generate

import (
	"fmt"
	"os/exec"
)

func (g generator) Build() error {
	fmt.Println("Building project...")
	cmd := exec.Command(
		"dotnet",
		"build",
		g.outPath,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("\nGenerated client '%s' to '%s'!\n", g.csModulePath, g.outPath)

	return nil
}

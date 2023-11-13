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
	return cmd.Run()
}

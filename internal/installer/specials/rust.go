package specials

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/sys"
	"github.com/yousefbustamiii/package-mate/internal/ui"
)

func rustCheck() (bool, string) {
	if _, err := exec.LookPath("rustc"); err == nil {
		out, _ := exec.Command("rustc", "--version").Output()
		return true, strings.TrimSpace(string(out))
	}
	return false, ""
}

func rustInstall() error {
	ui.Doing("Installing %s", "Rust")
	cmd := sys.ShellCommand("curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	ui.Done("Successfully installed %s", "Rust")
	return nil
}

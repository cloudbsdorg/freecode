//go:build !windows

package util

import (
	"os/exec"
)

func applyPlatformAttrs(cmd *exec.Cmd) {
}

package bluetooth

import (
	"context"
	"os/exec"
)

// onPath returns true if the named executable is on the $PATH. This will return
// false even if the error is not ErrNotFound, so issues could potentially arise
// in edge cases.
func onPath(executable string) bool {
	_, err := exec.LookPath(executable)
	return err == nil
}

func execcmd(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, cmd, args...).Output()
}

package execx

import (
	"bytes"
	"os/exec"
)

func Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func RunBash(cmd string) (string, error) {
	return Run("bash", "-c", cmd)
}

func RunPipes(pipes [][]string) (string, error) {
	var out bytes.Buffer
	for _, pipe := range pipes {
		cmd := exec.Command(pipe[0], pipe[1:]...)
		cmd.Stdin = &out
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		if err != nil {
			return out.String(), err
		}
	}
	return out.String(), nil
}

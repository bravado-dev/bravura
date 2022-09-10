package command

import (
	"io"
	"os/exec"
)

type RunOptions struct {
	RepoPath     string
	Environ      []string
	ExtraEnv     []string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
	ProcessGroup bool
}

type Command struct {
	*RunOptions
	rawCmd         *exec.Cmd
	closeAfterWait []io.Closer
}

func (c *Command) String() string {
	return c.rawCmd.String()
}

func (c *Command) Wait() error {
	err := c.rawCmd.Wait()
	_ = c.finalize()
	return err
}

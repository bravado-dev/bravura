//go:build linux

package command

func (c *Command) beforeInitialize() error {
	return nil
}

func (c *Command) delayInitialize() error {

	return nil
}

func (c *Command) finalize() error {

	return nil
}

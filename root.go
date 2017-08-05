package main

import (
	"github.com/spf13/cobra"
)

type GlobalOptions struct {
	Verbose  bool
	Insecure bool
}

type RootCommand struct {
	cmd  *cobra.Command
	opts *GlobalOptions
}

func NewRootCommand() *RootCommand {
	c := &RootCommand{
		cmd: &cobra.Command{
			Use:   "grpcurl",
			Short: "A handy and universal gRPC command line client",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Help()
			},
		},
		opts: &GlobalOptions{},
	}
	c.cmd.PersistentFlags().BoolVarP(&c.opts.Verbose, "verbose", "v", false, "verbose output")
	c.cmd.PersistentFlags().BoolVarP(&c.opts.Insecure, "insecure", "k", false, "with insecure")
	c.cmd.AddCommand(NewListServicesCommand(c.opts).Command())
	c.cmd.AddCommand(NewCallCommand(c.opts).Command())
	return c
}

func (c *RootCommand) Command() *cobra.Command {
	return c.cmd
}

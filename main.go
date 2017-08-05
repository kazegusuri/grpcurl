package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose  bool
	insecure bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "grpcurl",
		Short: "A handy and universal gRPC command line client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "with insecure")

	rootCmd.AddCommand(NewListServicesCommand().Command())
	rootCmd.AddCommand(NewCallCommand().Command())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package cmd

import (
	"encoding/json"
	"os"

	"github.com/jasonhancock/dyndns/version"
	"github.com/spf13/cobra"
)

func newVersionCmd(info version.Info) *cobra.Command {
	return &cobra.Command{
		Use:          "version",
		Short:        "Displays version information.",
		Long:         "Displays version information.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return json.NewEncoder(os.Stdout).Encode(info)
		},
	}
}

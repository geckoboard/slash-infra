package debug

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "debug utils",
	}

	cmd.AddCommand(eventbridgeCmd)

	return cmd
}

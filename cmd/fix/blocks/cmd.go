package blocks

import (
	"github.com/spf13/cobra"

	"github.com/forbole/juno/v2/cmd/parse"
)

// NewBlocksCmd returns the Cobra command that allows to fix all the things related to blocks
func NewBlocksCmd(parseConfig *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocks",
		Short: "Fix things related to blocks and transactions",
	}

	cmd.AddCommand(
		blocksCmd(parseConfig),
		heightsCmd(parseConfig),
	)

	return cmd
}

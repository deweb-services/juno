package blocks

import (
	"fmt"
	"strconv"

	"github.com/forbole/juno/v2/cmd/parse"

	"github.com/spf13/cobra"

	"github.com/forbole/juno/v2/parser"
)

// heightsCmd returns a Cobra command that allows to fix missing blocks in database
func heightsCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "heights [start height] [end height]",
		Short: "force process speficief block heights range",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parse.GetParsingContext(parseConfig)
			if err != nil {
				return err
			}

			workerCtx := parser.NewContext(parseCtx.EncodingConfig.Marshaler, nil, parseCtx.Node, parseCtx.Database, parseCtx.Logger, parseCtx.Modules)
			worker := parser.NewWorker(0, workerCtx)

			k, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("error while converting start height: %s", err)
			}
			endHeight, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("error while converting end height: %s", err)
			}

			fmt.Printf("Refetching missing blocks and transactions from height %d ... \n", k)
			for ; k <= endHeight; k++ {
				fmt.Println("Force processing height: ", k)
				err := worker.ForceProcess(k)
				if err != nil {
					return fmt.Errorf("error while re-fetching block %d: %s", k, err)
				}
			}

			return nil
		},
	}
}

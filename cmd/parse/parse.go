package parse

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/forbole/juno/v2/types/config"

	"github.com/go-co-op/gocron"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/juno/v2/modules"
	"github.com/forbole/juno/v2/parser"
	"github.com/forbole/juno/v2/types"

	"github.com/spf13/cobra"
)

var (
	waitGroup sync.WaitGroup
)

// ParseCmd returns the command that should be run when we want to start parsing a chain state.
func ParseCmd(cmdCfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:     "parse",
		Short:   "Start parsing the blockchain data",
		PreRunE: ReadConfig(cmdCfg),
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := GetParsingContext(cmdCfg)
			if err != nil {
				return err
			}

			// Run all the additional operations
			for _, module := range context.Modules {
				if module, ok := module.(modules.AdditionalOperationsModule); ok {
					err = module.RunAdditionalOperations()
					if err != nil {
						return err
					}
				}
			}

			return StartParsing(context)
		},
	}
}

// StartParsing represents the function that should be called when the parse command is executed
func StartParsing(ctx *Context) error {
	// Get the config
	cfg := config.Cfg.Parser

	// Start periodic operations
	scheduler := gocron.NewScheduler(time.UTC)
	for _, module := range ctx.Modules {
		if module, ok := module.(modules.PeriodicOperationsModule); ok {
			err := module.RegisterPeriodicOperations(scheduler)
			if err != nil {
				return err
			}
		}
	}
	scheduler.StartAsync()

	// Create a queue that will collect, aggregate, and export blocks and metadata
	exportQueue := types.NewQueue(25)

	// Create workers
	workerCtx := parser.NewContext(ctx.EncodingConfig.Marshaler, exportQueue, ctx.Node, ctx.Database, ctx.Logger, ctx.Modules)
	workers := make([]parser.Worker, cfg.Workers, cfg.Workers)
	for i := range workers {
		workers[i] = parser.NewWorker(i, workerCtx)
	}

	waitGroup.Add(1)

	// Run all the async operations
	for _, module := range ctx.Modules {
		if module, ok := module.(modules.AsyncOperationsModule); ok {
			go module.RunAsyncOperations()
		}
	}

	// Start each blocking worker in a go-routine where the worker consumes jobs
	// off of the export queue.
	for i, w := range workers {
		ctx.Logger.Debug("starting worker...", "number", i+1)
		go w.Start()
	}

	// Listen for and trap any OS signal to gracefully shutdown and exit
	trapSignal(ctx)

	if cfg.ParseGenesis {
		// Add the genesis to the queue if requested
		queueTask := types.NewQueueTask(0)
		exportQueue <- queueTask
	}

	go startNewBlockListener(exportQueue, ctx)

	// Block main process (signal capture will call WaitGroup's Done)
	waitGroup.Wait()
	return nil
}

// startNewBlockListener subscribes to new block events via the Tendermint RPCConfig
// and enqueues each new block height onto the provided queue. It blocks as new
// blocks are incoming.
func startNewBlockListener(exportQueue types.HeightQueue, ctx *Context) {
	eventCh, cancel, err := ctx.Node.SubscribeNewBlocks("juno-new-blocks-listener")
	defer cancel()

	if err != nil {
		panic(fmt.Errorf("failed to subscribe to new blocks: %s", err))
	}

	ctx.Logger.Info("listening for new block events...")

	for e := range eventCh {
		newBlock := e.Data.(tmtypes.EventDataNewBlock).Block
		height := newBlock.Header.Height

		ctx.Logger.Debug("enqueueing new block", "height", height)
		queueTask := types.NewQueueTask(height)
		exportQueue <- queueTask
	}
}

// trapSignal will listen for any OS signal and invoke Done on the main
// WaitGroup allowing the main process to gracefully exit.
func trapSignal(ctx *Context) {
	var sigCh = make(chan os.Signal)

	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		ctx.Logger.Info("caught signal; shutting down...", "signal", sig.String())
		defer ctx.Node.Stop()
		defer ctx.Database.Close()
		defer waitGroup.Done()
	}()
}

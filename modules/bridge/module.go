package bridge

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v2/database"
	"github.com/forbole/juno/v2/logging"
	"github.com/forbole/juno/v2/modules"
	"github.com/forbole/juno/v2/node"
	"github.com/forbole/juno/v2/types"
	"github.com/forbole/juno/v2/types/config"
)

var _ modules.Module = &Module{}

type ChainConfig struct {
	chain string
	token string
}

// Module represents the module allowing to store messages properly inside a dedicated table
type Module struct {
	cdc                    codec.Codec
	db                     database.Database
	node                   node.Node
	logger                 logging.Logger
	networkTokens          map[string]ChainConfig
	bridgeWalletAddress    string
	ConsensusModuleAddress string
}

func NewModule(cdc codec.Codec, db database.Database, bridgeConf config.BridgeConfig, node node.Node, logger logging.Logger) *Module {
	chainBridgeSettings := make(map[string]ChainConfig)
	for chain, networkConf := range bridgeConf.NetworksTokens {
		chainConf := ChainConfig{
			token: networkConf.Token,
			chain: chain,
		}
		chainBridgeSettings[networkConf.Token] = chainConf
	}
	return &Module{
		cdc:                    cdc,
		db:                     db,
		node:                   node,
		logger:                 logger,
		networkTokens:          chainBridgeSettings,
		bridgeWalletAddress:    bridgeConf.WalletAddress,
		ConsensusModuleAddress: bridgeConf.ConsensusModuleAddress,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "bridge_transactions"
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	return m.processMessage(msg, tx, m.cdc, m.db)
}

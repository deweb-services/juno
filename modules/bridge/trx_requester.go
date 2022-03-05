package bridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/forbole/juno/v2/types"
	"net/http"
)

type TransactionRequest struct {
	Chain   string `json:"chain"`
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

func (m *Module) processTokenTransferBridgeMsg(transferMsg *types.WasmTransferMsg, targetChain string) error {
	storedMapping, err := m.node.GetMappingToExternalAddress(transferMsg.Sender, targetChain)
	if err != nil || len(storedMapping.ExternalAddress) == 0 {
		return fmt.Errorf("error getting mapping value for address %s to chain %s",
			transferMsg.Sender, targetChain)
	}
	transferRequestMsg := &TransactionRequest{
		Chain:   targetChain,
		Address: storedMapping.ExternalAddress,
		Amount:  transferMsg.Amount,
	}
	_ = transferRequestMsg
	m.logger.Info(fmt.Sprintf("Sending request to transfer %s tokens in %s to %s",
		transferRequestMsg.Amount, transferRequestMsg.Chain, transferRequestMsg.Address))
	err = m.sendTransferRequest(transferRequestMsg)
	if err != nil {
		return fmt.Errorf("send transfer error: %w", err)
	}
	return nil
}

func (m *Module) sendTransferRequest(request *TransactionRequest) error {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshalling request: %w", err)
	}
	reqContent := bytes.NewBuffer(requestBody)
	resp, err := http.Post(m.ConsensusModuleAddress, "application/json", reqContent)
	if err != nil {
		return fmt.Errorf("error sending transafer request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("from consensus module received response code %d", resp.StatusCode)
	}
	return nil
}

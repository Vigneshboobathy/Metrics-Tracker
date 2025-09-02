	package rpc

	import (
		"encoding/json"
		"fmt"
		"metrics/models"
	)

	
	// GetAccountState
	func GetAccountState(node model.NodeInfo, address string) (map[string]interface{}, error) {
		req := model.RequestToRPC{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "xygle_getAccountState",
			Params:  map[string]interface{}{"address": address},
		}

		rpcResp, err := SendRequestToRPC(node.URL, req)
		if err != nil {
			return nil, err
		}

		var result map[string]interface{}
		if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal account state: %v", err)
		}
		return result, nil
	}

	// GetTransactionDetails
	func GetTransactionDetails(node model.NodeInfo, txID string) (model.TransactionResult, error) {
		req := model.RequestToRPC{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "xygle_getTransaction",
			Params:  map[string]interface{}{"id": txID},
		}

		rpcResp, err := SendRequestToRPC(node.URL, req)
		if err != nil {
			return model.TransactionResult{}, err
		}

		var result model.TransactionResult
		if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
			return model.TransactionResult{}, fmt.Errorf("failed to unmarshal transaction details: %v", err)
		}

		if result.ID == "" && result.TransactionID != "" {
			result.ID = result.TransactionID
		}

		return result, nil
	}

	// TransferFund
	func TransferFund(node model.NodeInfo, receiver string, value int, nonce int) (string, error) {
		req := model.RequestToRPC{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "xygle_transferFund",
			Params: map[string]interface{}{
				"receiver": receiver,
				"value":    value,
				"nonce":    nonce,
			},
		}
	
		rpcResp, err := SendRequestToRPC(node.URL, req)
		if err != nil {
			return "", err
		}
	
		var result struct {
			Status        string `json:"status"`
			TransactionID string `json:"transaction_id"`
		}
		if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
			return "", fmt.Errorf("failed to unmarshal transfer response: %v", err)
		}
	
		return result.TransactionID, nil
	}

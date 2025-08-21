package rpc

import (
	"encoding/json"
	"fmt"
	"metrics/models"
)

// GetBalance
func GetBalance(node model.NodeInfo, address string) (map[string]interface{}, error) {
	req := model.RequestToRPC{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "xygle_getBalance",
		Params:  map[string]interface{}{"address": address},
	}

	rpcResp, err := SendRequestToRPC(node.URL, req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal balance: %v", err)
	}
	return result, nil
}

// GetAllTransactions
func GetAllTransactions(node model.NodeInfo, page, pageSize int) (map[string]interface{}, error) {
	req := model.RequestToRPC{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "xygle_getAllTransactions",
		Params: map[string]interface{}{
			"page":     page,
			"pageSize": pageSize,
		},
	}

	rpcResp, err := SendRequestToRPC(node.URL, req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal all transactions: %v", err)
	}
	return result, nil
}

// GetAccountTransactions
func GetAccountTransactions(node model.NodeInfo, address string, page, pageSize int) (map[string]interface{}, error) {
	req := model.RequestToRPC{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "xygle_getAccountTransactions",
		Params: map[string]interface{}{
			"address":   address,
			"page":      page,
			"page_size": pageSize,
		},
	}

	rpcResp, err := SendRequestToRPC(node.URL, req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account transactions: %v", err)
	}
	return result, nil
}

// GetStats
func GetNodeStats(node model.NodeInfo) (map[string]interface{}, error) {
	req := model.RequestToRPC{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "xygle_getStats",
		Params:  map[string]interface{}{},
	}

	rpcResp, err := SendRequestToRPC(node.URL, req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal node stats: %v", err)
	}
	return result, nil
}

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
		return model.TransactionResult{}, err
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

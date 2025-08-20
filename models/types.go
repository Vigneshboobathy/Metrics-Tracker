package model

import "encoding/json"

type NodeInfo struct {
	NodeType string `json:"node_type"`
	URL      string `json:"url"`
	Address  string `json:"address"`
}

// RPCRequest represents a request to the RPC server
type RequestToRPC struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

// RPCResponse represents a response from the RPC server
type ResponseFromRPC struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data,omitempty"`
	} `json:"error"`
}

// GetTransactionResult represents a result of a transaction
type TransactionResult struct {
	ID        string `json:"id"`
	IsFinal   bool   `json:"is_final"`
	Timestamp int64  `json:"timestamp"`
}

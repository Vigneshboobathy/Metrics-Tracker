package model

import "encoding/json"

type NodeInfo struct {
	NodeType string `json:"node_type"`
	URL      string `json:"url"`
	Address  string `json:"address"`
}

type RequestToRPC struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

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

type TransactionResult struct {
	ID                 string `json:"id"`
	TransactionID      string `json:"transaction_id,omitempty"` 
	Sender             string `json:"sender,omitempty"`
	Receiver           string `json:"receiver,omitempty"`
	Value              int64  `json:"value,omitempty"`
	Nonce              uint64 `json:"nonce,omitempty"`
	Timestamp          int64  `json:"timestamp,omitempty"`            
	ExecutionStatus    string `json:"execution_status,omitempty"`    
	ExecutionResult    string `json:"execution_result,omitempty"`
	ExecutionTimestamp int64  `json:"execution_timestamp,omitempty"`  
	IsFinal            bool   `json:"is_final,omitempty"`
}

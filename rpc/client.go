package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
    "metrics/logger"
	"metrics/models"
)

func SendRequestToRPC(url string, req model.RequestToRPC) (model.ResponseFromRPC, error) {
	var rpcResp model.ResponseFromRPC

	body, err := json.Marshal(req)
	if err != nil {
		return rpcResp, fmt.Errorf("failed to marshal request: %v", err)
	}

	start := time.Now()
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		err = fmt.Errorf("failed to send request to %s: %v", url, err)
		logger.Error.Println(err)
		return rpcResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("RPC HTTP status %d: %s", resp.StatusCode, string(respBody))
		logger.Error.Println(err)
		return rpcResp, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response: %v", err)
		logger.Error.Println(err)
		return rpcResp, err
	}

	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		logger.Error.Printf("failed to unmarshal response: %v; body: %s", err, string(respBody))
		return rpcResp, fmt.Errorf("failed to unmarshal RPC response: %v", err)
	}

	if rpcResp.Error != nil {
		err := fmt.Errorf("RPC error: %s (data: %s)", rpcResp.Error.Message, rpcResp.Error.Data)
		logger.Error.Println(err)
		return rpcResp, err
	}

	if len(rpcResp.Result) == 0 || string(rpcResp.Result) == "null" {
		err := fmt.Errorf("RPC returned empty result for method %s", req.Method)
		logger.Error.Println(err)
		return rpcResp, err
	}

	// Log latency only for successful attempts
	latency := time.Since(start)
	LogLatency(req.Method, latency)

	return rpcResp, nil
}

// LogLatency logs latency for an RPC call
func LogLatency(method string, duration time.Duration) {
	logger.Metrics.Printf("Latency for %s = %v", method, duration)
}

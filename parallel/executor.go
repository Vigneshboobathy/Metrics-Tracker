package parallel

import (
	"fmt"
	"metrics/logger"
	"metrics/metricstracker"
	"metrics/models"
	"metrics/nonce"
	"metrics/rpc"
	"strings"
	"sync"
	"time"
)

type TransactionRequest struct {
	ID       int
	Receiver string
	Value    int
}

type TransactionResult struct {
	ID      int
	Nonce   int
	TxID    string
	Success bool
	Error   error
	Latency time.Duration
}

type ParallelExecutor struct {
	node         model.NodeInfo
	nonceManager *nonce.NonceManager
	tracker      *metricstracker.Tracker
	maxRetries   int
	baseBackoff  time.Duration
	nonceTimeout time.Duration
	workers      int
}

func NewParallelExecutor(node model.NodeInfo, workers int) (*ParallelExecutor, error) {
	nonceManager, err := nonce.NewNonceManager(node)
	if err != nil {
		return nil, fmt.Errorf("failed to create nonce manager: %v", err)
	}

	tracker := metricstracker.NewTracker(node)

	return &ParallelExecutor{
		node:         node,
		nonceManager: nonceManager,
		tracker:      tracker,
		maxRetries:   5,
		baseBackoff:  100 * time.Millisecond,
		nonceTimeout: 30 * time.Second,
		workers:      workers,
	}, nil
}

func (pe *ParallelExecutor) ExecuteTransactions(requests []TransactionRequest) ([]TransactionResult, error) {
	if len(requests) == 0 {
		return nil, fmt.Errorf("no transaction requests provided")
	}

	var results []TransactionResult
	var wg sync.WaitGroup
	resultMutex := sync.Mutex{}

	semaphore := make(chan struct{}, pe.workers)

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request TransactionRequest) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := pe.executeTransactionSequential(index, request)

			resultMutex.Lock()
			results = append(results, result)
			if result.Success {
				pe.tracker.MarkSubmitted(result.TxID, time.Now().Add(-result.Latency))
			}
			resultMutex.Unlock()
		}(i, req)

		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()

	// logger.Metrics.Printf("Completed transaction submission phase: %d successful, %d failed",
	// 	pe.countSuccessful(results), len(results)-pe.countSuccessful(results))

	return results, nil
}

func (pe *ParallelExecutor) executeTransactionSequential(_ int, req TransactionRequest) TransactionResult {
	startTime := time.Now()

	nonce := pe.nonceManager.AllocateNonce()

	logger.Metrics.Printf("Processing transaction %d with nonce %d", req.ID, nonce)

	for attempt := 1; attempt <= pe.maxRetries; attempt++ {
		txID, err := rpc.TransferFund(pe.node, req.Receiver, req.Value, nonce)
		if err != nil {
			if attempt < pe.maxRetries && pe.shouldRetry(err) {
				backoff := time.Duration(attempt) * pe.baseBackoff
				logger.Metrics.Printf("Transaction %d (nonce=%d) attempt %d failed, retrying in %v: %v",
					req.ID, nonce, attempt, backoff, err)
				time.Sleep(backoff)
				continue
			}

			pe.nonceManager.MarkFailed(nonce)
			return TransactionResult{
				ID:      req.ID,
				Nonce:   nonce,
				Success: false,
				Error:   fmt.Errorf("transaction failed after %d attempts: %v", attempt, err),
				Latency: time.Since(startTime),
			}
		}

		pe.nonceManager.MarkSubmitted(nonce, txID)
		// logger.Metrics.Printf("Transaction %d submitted (nonce=%d) txID=%s", req.ID, nonce, txID)

		return TransactionResult{
			ID:      req.ID,
			Nonce:   nonce,
			TxID:    txID,
			Success: true,
			Latency: time.Since(startTime),
		}
	}

	pe.nonceManager.MarkFailed(nonce)
	return TransactionResult{
		ID:      req.ID,
		Nonce:   nonce,
		Success: false,
		Error:   fmt.Errorf("unexpected execution path"),
		Latency: time.Since(startTime),
	}
}

func (pe *ParallelExecutor) shouldRetry(err error) bool {
	errStr := strings.ToLower(err.Error())

	retryableErrors := []string{
		"connection refused",
		"timeout",
		"temporary",
		"network",
		"502",
		"503",
		"504", 
	}

	for _, retryable := range retryableErrors {
		if strings.Contains(errStr, retryable) {
			return true
		}
	}

	nonRetryableErrors := []string{
		"nonce too low",
		"nonce mismatch",
		"invalid nonce",
		"insufficient funds",
		"invalid signature",
	}

	for _, nonRetryable := range nonRetryableErrors {
		if strings.Contains(errStr, nonRetryable) {
			return false
		}
	}
	return true
}

func (pe *ParallelExecutor) WaitForCompletion() (int, int) {
	logger.Metrics.Printf("Waiting for transaction completion...")

	stopMonitoring := make(chan bool)
	go pe.monitorExecutions(stopMonitoring)

	executed, finalized := pe.tracker.WaitAndCollect()

	close(stopMonitoring)

	// logger.Metrics.Printf("Transaction completion: Executed=%d, Finalized=%d", executed, finalized)
	return executed, finalized
}

func (pe *ParallelExecutor) monitorExecutions(stop <-chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			pe.updateNonceStates()
		}
	}
}

func (pe *ParallelExecutor) updateNonceStates() {
	states := pe.nonceManager.GetAllStates()

	for nonce, state := range states {
		state.Mutex.RLock()
		submitted := state.Submitted
		executed := state.Executed
		txID := state.TxID
		state.Mutex.RUnlock()

		if submitted && !executed && txID != "" {
			detail, err := rpc.GetTransactionDetails(pe.node, txID)
			if err != nil {
				continue
			}

			if detail.ExecutionStatus == "SUCCESS" {
				pe.nonceManager.MarkExecuted(nonce)
			}
		}
	}
}

func (pe *ParallelExecutor) GetTracker() *metricstracker.Tracker {
	return pe.tracker
}


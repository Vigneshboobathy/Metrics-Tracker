package main

import (
	"fmt"
	"metrics/config"
	"metrics/logger"
	"metrics/models"
	"metrics/parallel"
	"metrics/rpc"
)

func main() {
	logger.Init()

	// Properly Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	validatorNodes := model.NodeInfo{
		NodeType: cfg.Node.Type,
		URL:      cfg.Node.URL,
		Address:  cfg.Node.Address,
	}
	receiver := cfg.Receiver
	value := 1
	numTx := 10
	workers := 1

	// Get transaction details
	txDetail, _ := rpc.GetTransactionDetails(validatorNodes, "2b3210cc4c19d169765ddf3d4a01472356dc0263bf926ab2df8f309e27091e4a")
	logger.Metrics.Printf("txDetail: %+v", txDetail)

	// Create parallel executor
	executor, err := parallel.NewParallelExecutor(validatorNodes, workers)
	if err != nil {
		logger.Metrics.Printf("Failed to create parallel executor: %v", err)
		return
	}

	// Prepare transaction requests
	var requests []parallel.TransactionRequest
	for i := 1; i <= numTx; i++ {
		requests = append(requests, parallel.TransactionRequest{
			ID:       i,
			Receiver: receiver,
			Value:    value,
		})
	}

	logger.Metrics.Printf("Starting sequential execution of %d transactions with %d workers", numTx, workers)

	// Execute transactions with proper nonce coordination
	results, err := executor.ExecuteTransactions(requests)
	if err != nil {
		logger.Metrics.Printf("Failed to execute transactions: %v", err)
		return
	}

	// Log submission results
	successful := 0
	failed := 0
	for _, result := range results {
		if result.Success {
			successful++
			logger.Metrics.Printf("Transaction %d submitted successfully (nonce=%d, txID=%s, latency=%.3fs)",
				result.ID, result.Nonce, result.TxID, result.Latency.Seconds())
		} else {
			failed++
			logger.Metrics.Printf("Transaction %d failed (nonce=%d): %v", result.ID, result.Nonce, result.Error)
		}
	}

	logger.Metrics.Printf("Submission phase completed: %d successful, %d failed", successful, failed)

	// Wait for execution and finalization
	executed, finalized := executor.WaitForCompletion()
	logger.Metrics.Printf("Execution phase completed: Executed=%d, Finalized=%d", executed, finalized)

	// Generate summary
	tracker := executor.GetTracker()
	sum, err := tracker.Summarize()
	if err != nil {
		logger.Metrics.Printf("Summary note: %v", err)
	}

	// Per-transaction metrics
	for txID, lat := range sum.LatencySeconds {
		logger.Metrics.Printf("Tx %s latency=%.2fs", txID, lat)
	}
	for txID, fin := range sum.TimeToFinalSeconds {
		logger.Metrics.Printf("Tx %s time_to_final=%.2fs", txID, fin)
	}

	// Summary metrics
	logger.Metrics.Printf("PERFORMANCE SUMMARY")
	logger.Metrics.Printf("Total submitted: %d, Successful: %d, Failed: %d", len(results), successful, failed)
	logger.Metrics.Printf("Executed: %d, Finalized: %d", sum.ExecutedCount, sum.FinalizedCount)
	logger.Metrics.Printf("Average latency: %.2fs over %d executed txs", sum.AvgLatencySeconds, sum.ExecutedCount)
	if sum.FinalizedCount > 0 {
		logger.Metrics.Printf("Average time-to-finality: %.2fs over %d finalized txs", sum.AvgTimeToFinalSeconds, sum.FinalizedCount)
	} else {
		logger.Metrics.Printf("No txs reached finality within the timeout window")
	}
	logger.Metrics.Printf("Estimated TPS: %.2f", sum.TPS)
}

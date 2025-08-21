package main

import (
	"metrics/logger"
	"metrics/metrics"
	"metrics/models"
	"metrics/rpc"
	"time"
)

var (
	validatorNodes = model.NodeInfo{
		NodeType: "validator",
		URL:      "http://103.189.172.20:8545",
		Address:  "6fe927e4f99777b8dc152963306a0b79e20dfdcb",
	}
)

func main() {
	logger.Init()
	metrics.StartTPSLogger()

	// Transfer funds
	receiver := "34187346aa86c5720a4d7b477e60060a8fadeb14"
	value := 10
	numTx := 5

	// Get transaction details
	txDetail, _ := rpc.GetTransactionDetails(validatorNodes, "2b3210cc4c19d169765ddf3d4a01472356dc0263bf926ab2df8f309e27091e4a")
	logger.Metrics.Printf("txDetail: %+v", txDetail)

	// Get balance
	balance, _ := rpc.GetBalance(validatorNodes, "38fa5e14098da08f0f241df3539cc6ee3a6f29d3")
	logger.Metrics.Printf("Balance: %+v", balance)

	// Get all transactions
	allTx, _ := rpc.GetAllTransactions(validatorNodes, 1, 10)
	logger.Metrics.Printf("fetched all transcations total =%v page=%v", allTx["total"], allTx["page"])

	// Get account transactions
	accTx, _ := rpc.GetAccountTransactions(validatorNodes, "6fe927e4f99777b8dc152963306a0b79e20dfdcb", 1, 1)
	logger.Metrics.Printf("Account Transactions: %+v", accTx)

	// Get node stats
	stats, _ := rpc.GetNodeStats(validatorNodes)
	logger.Metrics.Printf("Node Stats: %+v", stats)

	// Get account state
	state, _ := rpc.GetAccountState(validatorNodes, "d6e1d9c4bba0983991ce8d3fb8d0345009654387")
	logger.Metrics.Printf("Account state: %+v", state)

	// Get account state for current nonce
	state, err := rpc.GetAccountState(validatorNodes, validatorNodes.Address)
	if err != nil {
		logger.Metrics.Printf("Failed to fetch account state: %v", err)
		return
	}
	currentNonce := int(state["nonce"].(float64))

	txsubmission := make(map[string]time.Time)
	startBatch := time.Now()
	// Submit multiple txs
	for i := 1; i <= numTx; i++ {
		nonce := currentNonce + i
		txID, err := rpc.TransferFund(validatorNodes, receiver, value, nonce)
		if err != nil {
			logger.Metrics.Printf("Transaction %d failed (nonce=%d): %v", i, nonce, err)
		} else {
			logger.Metrics.Printf("Transaction %d submitted successfully (nonce=%d): txID=%s", i, nonce, txID)
			txsubmission[txID] = time.Now()
		}
	}

	// TPS for this batch
	elapsed := time.Since(startBatch).Seconds()
	tps := float64(len(txsubmission)) / elapsed
	logger.Metrics.Printf("Batch TPS = %.2f (submitted %d txs in %.2fs)", tps, len(txsubmission), elapsed)

	// Finality Check
	var totalFinality time.Duration
	var firstTx, lastTx time.Time

	for _, submittedAt := range txsubmission {
		if firstTx.IsZero() || submittedAt.Before(firstTx) {
			firstTx = submittedAt
		}
		if lastTx.IsZero() || submittedAt.After(lastTx) {
			lastTx = submittedAt
		}
		// measure how long it took this tx to be submitted
		finality := submittedAt.Sub(startBatch)
		totalFinality += finality
	}

	logger.Metrics.Printf("Total Finality (sum of all txs) = %v", totalFinality)

}

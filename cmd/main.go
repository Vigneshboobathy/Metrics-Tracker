package main

import (
	"metrics/logger"
	"metrics/metrics"
	"metrics/models"
	"metrics/rpc"
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

	// Get transaction details
	txDetail, _ := rpc.GetTransactionDetails(validatorNodes, "1")
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

	// Transfer Funds
	txID, err := rpc.TransferFund(
		validatorNodes,
		"d6e1d9c4bba0983991ce8d3fb8d0345009654387",
		100,
		16,
	)
	if err != nil || txID == "" {
		logger.Metrics.Printf("No transaction has been submitted (error: %v)", err)
	} else {
		logger.Metrics.Printf("Submitted tx: %s", txID)
	}
}

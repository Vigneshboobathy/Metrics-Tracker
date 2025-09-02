package nonce

import (
	"fmt"
	"metrics/logger"
	"metrics/models"
	"metrics/rpc"
	"sync"
	"time"
)

type NonceManager struct {
	node         model.NodeInfo
	currentNonce int
	mutex        sync.Mutex
	nonceStates map[int]*NonceState
	statesMutex sync.RWMutex
}

type NonceState struct {
	Nonce       int
	TxID        string
	Submitted   bool
	Executed    bool
	Failed      bool
	SubmittedAt time.Time
	ExecutedAt  time.Time
	Mutex       sync.RWMutex 
}

func NewNonceManager(node model.NodeInfo) (*NonceManager, error) {
	state, err := rpc.GetAccountState(node, node.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get account state: %v", err)
	}

	currentNonce := int(state["nonce"].(float64))

	return &NonceManager{
		node:         node,
		currentNonce: currentNonce,
		nonceStates:  make(map[int]*NonceState),
	}, nil
}

func (nm *NonceManager) AllocateNonce() int {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	nm.currentNonce++
	nonce := nm.currentNonce

	nm.statesMutex.Lock()
	nm.nonceStates[nonce] = &NonceState{
		Nonce: nonce,
	}
	nm.statesMutex.Unlock()

	logger.Metrics.Printf("Allocated nonce %d", nonce)
	return nonce
}

func (nm *NonceManager) MarkSubmitted(nonce int, txID string) {
	nm.statesMutex.RLock()
	state, exists := nm.nonceStates[nonce]
	nm.statesMutex.RUnlock()

	if !exists {
		logger.Error.Printf("Nonce %d not found in states", nonce)
		return
	}

	state.Mutex.Lock()
	state.TxID = txID
	state.Submitted = true
	state.SubmittedAt = time.Now()
	state.Mutex.Unlock()
}

func (nm *NonceManager) MarkExecuted(nonce int) {
	nm.statesMutex.RLock()
	state, exists := nm.nonceStates[nonce]
	nm.statesMutex.RUnlock()

	if !exists {
		logger.Error.Printf("Nonce %d not found in states", nonce)
		return
	}

	state.Mutex.Lock()
	state.Executed = true
	state.ExecutedAt = time.Now()
	state.Mutex.Unlock()

	// logger.Metrics.Printf("Marked nonce %d as executed", nonce)
}

func (nm *NonceManager) MarkFailed(nonce int) {
	nm.statesMutex.RLock()
	state, exists := nm.nonceStates[nonce]
	nm.statesMutex.RUnlock()

	if !exists {
		logger.Error.Printf("Nonce %d not found in states", nonce)
		return
	}

	state.Mutex.Lock()
	state.Failed = true
	state.Mutex.Unlock()

	logger.Metrics.Printf("Marked nonce %d as failed", nonce)
}

func (nm *NonceManager) GetAllStates() map[int]*NonceState {
	nm.statesMutex.RLock()
	defer nm.statesMutex.RUnlock()
	result := make(map[int]*NonceState)
	for k, v := range nm.nonceStates {
		result[k] = v
	}
	return result
}



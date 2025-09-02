package metricstracker

import (
	"errors"
	"metrics/logger"
	"metrics/models"
	"metrics/rpc"
	"sort"
	"strings"
	"time"
)

type txTimes struct {
	submitted  time.Time
	executed   time.Time
	finalized  time.Time
	execUnix   int64
	execStatus string
}

type Tracker struct {
	node      model.NodeInfo
	times     map[string]*txTimes
	pollEvery time.Duration
	timeout   time.Duration
}

type Summary struct {
	LatencySeconds        map[string]float64
	TimeToFinalSeconds    map[string]float64
	ExecUnixTimestamps    map[string]int64
	AvgLatencySeconds     float64
	AvgTimeToFinalSeconds float64
	TPS                   float64
	ExecutedCount         int
	FinalizedCount        int
}

func NewTracker(node model.NodeInfo) *Tracker {
	return &Tracker{
		node:      node,
		times:     make(map[string]*txTimes),
		pollEvery: 2 * time.Second,
		timeout:   5 * time.Minute,
	}
}

func (t *Tracker) MarkSubmitted(txID string, at time.Time) {
	t.times[txID] = &txTimes{submitted: at}
}

func (t *Tracker) WaitAndCollect() (int, int) {
	deadline := time.Now().Add(t.timeout)
	executed := 0
	final := 0

	pending := func() []string {
		var ids []string
		for id, tt := range t.times {
			if tt.executed.IsZero() || tt.finalized.IsZero() {
				ids = append(ids, id)
			}
		}
		return ids
	}

	for {
		if len(pending()) == 0 || time.Now().After(deadline) {
			break
		}

		for _, txID := range pending() {
			detail, err := rpc.GetTransactionDetails(t.node, txID)
			if err != nil {
				if !strings.Contains(strings.ToLower(err.Error()), "not found") {
					logger.Error.Printf("poll %s error: %v", txID, err)
				}
				continue
			}

			tt := t.times[txID]
			if detail.ExecutionStatus == "SUCCESS" && tt.executed.IsZero() {
				tt.executed = time.Now()
				tt.execUnix = detail.ExecutionTimestamp
				tt.execStatus = detail.ExecutionStatus
				executed++
				logger.Metrics.Printf("Tx %s executed (status=%s)", txID, detail.ExecutionStatus)
			}
			if detail.IsFinal && tt.finalized.IsZero() {
				tt.finalized = time.Now()
				final++
				logger.Metrics.Printf("Tx %s is final =%t)", txID, detail.IsFinal)
			}
		}

		time.Sleep(t.pollEvery)
	}

	return executed, final
}

func (t *Tracker) Summarize() (Summary, error) {
	s := Summary{
		LatencySeconds:     map[string]float64{},
		TimeToFinalSeconds: map[string]float64{},
		ExecUnixTimestamps: map[string]int64{},
	}
	var latVals []float64
	var finVals []float64

	var execTs []int64
	for id, tt := range t.times {
		if !tt.executed.IsZero() && !tt.submitted.IsZero() {
			lat := tt.executed.Sub(tt.submitted).Seconds()
			s.LatencySeconds[id] = lat
			latVals = append(latVals, lat)
			s.ExecUnixTimestamps[id] = tt.execUnix
			if tt.execUnix > 0 {
				execTs = append(execTs, tt.execUnix)
			}
			s.ExecutedCount++
		}
		if !tt.finalized.IsZero() && !tt.submitted.IsZero() {
			final := tt.finalized.Sub(tt.submitted).Seconds()
			s.TimeToFinalSeconds[id] = final
			finVals = append(finVals, final)
			s.FinalizedCount++
		}
	}

	if len(latVals) > 0 {
		var sum float64
		for _, v := range latVals {
			sum += v
		}
		s.AvgLatencySeconds = sum / float64(len(latVals))
	}
	if len(finVals) > 0 {
		var sum float64
		for _, v := range finVals {
			sum += v
		}
		s.AvgTimeToFinalSeconds = sum / float64(len(finVals))
	}

	if len(execTs) == 0 {
		return s, errors.New("no executed transactions to compute TPS")
	}
	sort.Slice(execTs, func(i, j int) bool { return execTs[i] < execTs[j] })
	minTs := execTs[0]
	maxTs := execTs[len(execTs)-1]
	span := maxTs - minTs
	if span <= 0 {
		s.TPS = float64(len(execTs))
	} else {
		s.TPS = float64(len(execTs)) / float64(span)
	}

	return s, nil
}

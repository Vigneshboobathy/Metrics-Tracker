### MetricsTracker

A Go service to submit transactions in parallel and record key client-side DAG metrics:

- **Metrics logging**: TPS, RPC latency, and time-to-finality written to `metrics.log`.
- **Parallel executor**: Coordinated nonces for safe concurrent submissions.
- **RPC helpers**: Balance lookup, transactions listing, account activity, node stats, and transaction details.

### Project structure

- `cmd/main.go`: Example program wiring logging, RPC calls, and the parallel executor
- `logger/`: Console and file logger; writes metrics to `metrics.log`
- `metricstracker/`: Aggregates timings and computes summary metrics (latency, time-to-finality, TPS)
- `models/`: Shared request/response and type definitions
- `parallel/`: Parallel transaction executor with nonce coordination and completion tracking
- `rpc/`: HTTP JSON-RPC client and high-level helpers
- `metrics.log`: Metrics output file created at runtime

### Requirements

- Go 1.23+

### Setup

### Clone the repository
https://github.com/Vigneshboobathy/Metrics-Tracker.git

### Open a terminal and move into the module directory:

```bash
cd Metrics-Tracker
```

```bash
go mod tidy
```

### Run

```bash
go run cmd/main.go
```
```go
validatorNodes := model.NodeInfo{
    NodeType: "validator",
    URL:      "http://<rpc-host>:<port>",
    Address:  "<sender-address>",
}

receiver := "<receiver-address>"
value := 1
numTx := 25
workers := 3
```

### What it does at runtime

- Initializes logging to console and `metrics.log`.
- Performs sample RPC calls: transaction details, balance, transactions list, account activity, and node stats.
- Submits `numTx` transactions via the parallel executor, coordinating nonces across `workers`.
- Waits for execution and finality, polling transaction status periodically.
- Produces a performance summary including per-transaction latencies and aggregate metrics.

### Output

- Console logs and `metrics.log` will include lines like:
  - `Transaction <id> submitted successfully (nonce=<n>, txID=<hash>, latency=<s>)`
  - `Tx <hash> executed (status=SUCCESS)`
  - `Tx <hash> is final =true`
  - `PERFORMANCE SUMMARY` with counts, averages, and `Estimated TPS`

### Metrics explained

- **Latency (s)**: Time from client submission to observed execution.
- **Time-to-finality (s)**: Time from submission until the transaction is final.
- **TPS**: Derived from execution timestamps over the observed window.


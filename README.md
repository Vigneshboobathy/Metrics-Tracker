### MetricsTracker

A minimal Go service for tracking and logging client-side blockchain metrics (TPS, RPC latency, and time-to-finality) while providing convenience RPC helpers for common operations.

- **Metrics logging**: Writes TPS, RPC latencies, and finality times to `metrics.log`.
- **RPC helpers**: Balance lookup, transactions listing, account state, node stats, transaction details, and simple fund 

### Project structure

- `MetricsTracker/cmd/main.go`: example program wiring logging + RPC calls
- `MetricsTracker/logger/`: stdout/stderr + metrics file logger
- `MetricsTracker/metrics/`: TPS counter and metrics emitters
- `MetricsTracker/models/`: request/response and shared types
- `MetricsTracker/rpc/`: HTTP JSON-RPC client and high-level helpers
- `MetricsTracker/metrics.log`: created at runtime by the app

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

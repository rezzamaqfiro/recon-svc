# 🏦 Reconciliation Service

A high-integrity reconciliation engine built in Go to identify discrepancy between internal ledger transactions and external bank statements.

---

## ⚙️ Core Logic Explanation

### 1. Data Normalization
The service treats all currency as `int64` (IDR). By avoiding floating-point math, it eliminate rounding errors.

### 2. Two-Stage Matching Process
* **Temporal Anchoring:** The engine groups system transactions by date into a Hash Map ($O(n)$ complexity). It looks for matches within a configurable `windowDays` (e.g., +/- 1 day).
* **Heuristic:** If multiple transactions exist within the time window, the engine calculates the absolute difference between amounts and selects the "closest" candidate.

### 3. Consumption Logic
Once a record is matched, it is "consumed" (removed from the candidate pool). This prevents a single bank record from being matched against multiple system records, ensuring audit-ready totals.


---

## 🛠 Project Structure

```text
recon-svc/
├── main.go          # CLI Entrypoint & Orchestration
├── engine.go        # Recon Process
├── parser.go        # CSV Ingestion
├── models.go        # Domain Data Structures
├── engine_test.go   # Table-Driven Unit Tests
└── cmd/
    └── generator/
        └── gen_data.go      # The 10k data generation tool
```

## 🚀 How to Use
### 1. Preparation
Ensure you have provide your system.csv and bank.csv files in the root directory.

OR

You can use CLI and execute below command to generate data

```bash
go run cmd/generator/gen_data.go
```

We can put our desire row data by put additional flag parameter `default: 10000 rows`

```bash
go run cmd/generator/gen_data.go -rows=100
```

### 2. Running Detailed Tests
This project uses Table-Driven tests to cover edge cases like settlement lag and amount typos.

```bash
go test -v .
```

### 3. Running the Program
Use the CLI flags to specify your files.

```bash
go run . -system=system.csv -bank=bank_a.csv,bank_b.csv,bank_c.csv -start=2025-01-01 -end=2025-01-31
```

You can also extend the time window by adding configurable parameter `default: 1 days`

```bash
go run . -system=system.csv -bank=bank_a.csv,bank_b.csv,bank_c.csv -start=2025-01-01 -end=2025-01-31 -window=2
```

### 4. Build the Program
By executing below script you can build the progam into executable files.

```bash
go build -o recon-tool ./recon-tool -system=system.csv -bank=bank_a.csv,bank_b.csv,bank_c.csv -start=2025-01-01 -end=2025-01-31
```


---

**Author:** Rezza Maqfiro — 2026
# Wallet Service

## Overview

This project implements a simple **REST API** for managing digital wallets. It provides functionality for creating wallets, checking their balances, making deposits, and withdrawals. The backend uses **PostgreSQL** as the primary database and utilizes **Go**, specifically the **pgx** library, for efficient communication with the database. Additionally, the **Gin** web framework is employed to handle HTTP requests effectively.

## Features

- Creation of new wallets.
- Retrieval of wallet balances.
- Depositing funds into wallets.
- Withdrawing funds from wallets.
- Logging and basic error handling mechanisms.

## Quick Start Guide

### Prerequisites:
- Go (≥ 1.24.6);
- Postgresql;
- Goose (optional);
- Docker (optional).

### Running Locally:

- **Step 1**: Create a `config.env` file with environment variables, for example:

```bash
BIND_IP=127.0.0.1
LISTEN_PORT=8888
PSQL_HOST=127.0.0.1
PSQL_PORT=5432
PSQL_NAME=simple_wallet
PSQL_USER=user
PSQL_PASSWORD=password
```

- **Step 2**: Install `goose` migration tool (optional):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

- **Step 3**: Apply database migrations (optional):

```bash
goose -dir=migrations postgres \
"host=your_db_host port=your_db_port dbname=db_name user=your_username_here password=your_password_here sslmode=disable" up
```

- **Step 4**: Build and run the server:

```bash
make build
make run
```

---

### Running with Docker:

- **Step 1**: Create a `config.env` file with environment variables, for example:

```bash
BIND_IP=127.0.0.1
LISTEN_PORT=8888
PSQL_HOST=psql-db
PSQL_PORT=5432
PSQL_NAME=simple_wallet
PSQL_USER=user
PSQL_PASSWORD=password
```

- **Step 2**: Start the containerized application:

```bash
make docker-compose-up
```

---

## Usage

The following API endpoints are available:

| Method | Endpoint                   | Description                               |
|--------|----------------------------|-------------------------------------------|
| GET    | `/api/v1/wallets/{uuid}`   | Retrieve wallet balance                    |
| POST   | `/api/v1/wallet`           | Perform a transaction                     |
| POST   | `/api/v1/wallets`          | Create a new wallet                       |

### Request Body for Transactions (`POST /api/v1/wallet`)

When performing a transaction (deposit or withdrawal), provide the following body in JSON format:

```
{
  "valletId": "<Wallet UUID>",
  "operationType": "DEPOSIT" or "WITHDRAW",
  "amount": amount
}
```

Where:

- `valletId`: Unique identifier of the wallet.
- `operationType`: Type of transaction ("DEPOSIT" or "WITHDRAW").
- `amount`: Positive value representing the amount being deposited or withdrawn.

Example valid request bodies:

```json
{
  "valletId": "c5a72fdd-f1d8-47b2-b461-c132429120bb",
  "operationType": "DEPOSIT",
  "amount": 100.00
}
```

```json
{
  "valletId": "c5a72fdd-f1d8-47b2-b461-c132429120bb",
  "operationType": "WITHDRAW",
  "amount": 50.00
}
```

Чтобы добавить информацию о полном покрытии тестами хэндлеров и сервисов в ваш проект, можно создать отдельный раздел документации под названием **Testing**. Вот пример того, куда и как можно вставить этот раздел в существующую структуру README.md или документа проекта:

---

## Testing

All handler and service layers have been fully covered by unit tests using Go's built-in testing package. These tests ensure that every critical path through the code has at least one test case associated with it, providing complete coverage for all handlers and services within this project.

To execute these tests locally, use command:

```bash
make test_coverage
```

Additionally, you can generate detailed test coverage reports by running:

```bash
make test_report
```

These commands will output an HTML report detailing which parts of the code were executed during testing, ensuring confidence in both reliability and robustness of our implementation.

---

## Load Testing Results

Load testing was conducted on the system using Vegeta, a high-performance HTTP load testing tool. Below are the results achieved under heavy load conditions:

```shell
Requests      [total, rate, throughput]         30000, 1000.01, 1000.00
Duration      [total, attack, wait]             30s, 30s, 440.763µs
Latencies     [min, mean, 50, 90, 95, 99, max]  286.206µs, 1.24ms, 1.125ms, 2.029ms, 2.699ms, 4.453ms, 88.934ms
Bytes In      [total, mean]                     2742891, 91.43
Bytes Out     [total, mean]                     2170000, 72.33
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:30000
Error Set: Empty
```

### Key Observations:
- All requests were successfully handled without any failures.
- Average latency remained consistently low (~1 ms per request) despite heavy load.
- System maintained steady throughput at about 1000 requests per second.

## Reproducing Load Tests

Follow these steps to replicate similar benchmarks yourself:

### Step 1: Install Vegeta

Install Vegeta globally:

```bash
go install github.com/tsenart/vegeta@latest
```

Or manage it via dependency tools within your project setup.

### Step 2: Automatically Generate Test Wallet

A test wallet is automatically created when applying database migrations. This wallet is utilized for load testing purposes.

### Step 3: Launch Load Test

Run the load test using the predefined Makefile target:

```bash
make vegeta_test
```

By default, this command executes Vegeta with preset configurations (request count, duration, etc.).

### Step 4: Review Detailed Report

Generate and review the full test report with:

```bash
make vegeta_report
```

The report includes response time distribution, success ratios, and additional useful metrics.

---

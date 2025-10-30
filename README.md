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

```json
{
  "valletId": "<Wallet UUID>",
  "operationType": "DEPOSIT" OR "WITHDRAW",
  "amount": <Amount Value>
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

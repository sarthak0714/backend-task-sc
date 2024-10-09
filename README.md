# Portfolio Tracking API

This project implements a portfolio tracking API that allows adding, deleting, and updating trades, as well as performing basic return calculations for a single user's portfolio.

[Deployed Link](https://backend-task-sc-production.up.railway.app/)

## Overview

- Add, update, and remove trades for securities
- Fetch all trades for a user
- Fetch portfolio summary
- Calculate cumulative returns
- Input validation to ensure portfolio integrity

## Tech Stack

- Go (Golang) 1.22.2
- Echo framework for HTTP routing
- GORM as ORM
- PostgreSQL for database (Using neondb)
- Swagger for API documentation

## API Documentation

API documentation is available via Swagger UI. After starting the server, visit:

```
http://localhost:8080/
```

## Getting Started

### Prerequisites

- Go 1.22.2 or later
- PostgreSQL
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/sarthak0714/backend-task-sc.git
   cd backend-task-sc
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up your environment variables in a `.env` file:
   ```
   DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
   PORT=8080
   ```

4. Build the project:
   ```bash
   make build
   ```

## Usage

To run the server: 
```bash
make run
```
Or without Make:
```golang
go run cmd/main.go
```


## API Endpoints

- `GET /status`: Check API status
- `POST /trades`: Add a new trade
- `PUT /trades/:id`: Update an existing trade
- `DELETE /trades/:id`: Remove a trade
- `GET /trades/:userId`: Fetch all trades for a user
- `GET /portfolio/:userId`: Fetch user's portfolio
- `GET /returns`: Calculate cumulative returns

For detailed request/response formats, please refer to the Swagger documentation.


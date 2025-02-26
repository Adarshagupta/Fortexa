# Fortexa
A Distributed, Fault-Tolerant, Event-Driven Transaction Processing System

## Overview

Fortexa is a scalable payment processing system designed for modern fintech applications. This MVP implements the core components of a distributed payment infrastructure including:

- Multi-payment method support (Credit Cards, UPI, Bank Transfers)
- Event-driven architecture using Kafka
- Real-time transaction processing
- Basic fraud detection
- Settlement processing
- Merchant dashboard

## System Architecture

```
        +--------------------+
        |  Merchant Dashboard|
        +---------+----------+
                  |
                  v
        +--------------------+
        |  API Gateway       |
        +--------------------+
                  |
        +------------------------------+
        |  Payment Processing Engine    |
        |  (Go Microservices)           |
        +------------------------------+
                  |
        +------------------------------+
        |  Message Queue (Kafka)        |
        +------------------------------+
                  |
  +-----------------------+------------------+-----------------+
  | Settlement Engine     | Fraud Detection  | Logging & Audit |
  |-----------------------|------------------|-----------------|
  | Handles fund flows    | Basic fraud      | Stores events   |
  | & reconciliations     | detection        | for compliance  |
  +----------------------------------------------------------+
```

## Tech Stack

- **API Gateway**: Go with Gin framework
- **Backend Services**: Go microservices 
- **Message Queue**: Apache Kafka
- **Database**: PostgreSQL
- **Merchant Dashboard**: React
- **Containerization**: Docker & Docker Compose

## Project Structure

```
fortexa/
├── api-gateway/          # API Gateway service
├── payment-engine/       # Core payment processing service
├── settlement-engine/    # Handles settlements and reconciliation
├── fraud-detection/      # Basic fraud detection service
├── merchant-dashboard/   # React-based merchant UI
└── infrastructure/       # Infrastructure components
    ├── db/               # Database migrations and schemas
    └── kafka/            # Kafka configuration
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.20+
- Node.js 18+

### Setup & Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/fortexa.git
   cd fortexa
   ```

2. Start the infrastructure services:
   ```
   docker-compose up -d
   ```

3. Start the backend services:
   ```
   cd api-gateway && go run main.go
   ```
   (Repeat for other services)

4. Start the merchant dashboard:
   ```
   cd merchant-dashboard
   npm install
   npm start
   ```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/payments/initiate` | POST | Create a new payment |
| `/api/v1/payments/:id` | GET | Get payment status |
| `/api/v1/refunds` | POST | Issue a refund |
| `/api/v1/merchants/onboard` | POST | Onboard a new merchant |
| `/api/v1/webhooks/register` | POST | Register merchant webhook |

## MVP Features

- [x] Payment initiation and processing
- [x] Basic fraud checks
- [x] Transaction state management
- [x] Merchant onboarding
- [x] Settlement processing
- [x] Transaction history and reporting
- [x] Webhook notifications

## License

MIT

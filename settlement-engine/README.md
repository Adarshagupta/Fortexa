# Settlement Engine Service

This service is responsible for processing captured payments into settlements for merchants. It calculates fees, applies taxes, and manages settlement cycles.

## Environment Setup

The service uses environment variables for configuration. Create a `.env` file in the root of the project with the following variables:

```
# Application settings
APP_NAME=fortexa
APP_ENV=development
SHUTDOWN_TIMEOUT=10

# Database settings
DB_HOST=ep-rapid-dew-a1roz96h-pooler.ap-southeast-1.aws.neon.tech
DB_PORT=5432
DB_USER=neondb_owner
DB_PASSWORD=npg_QMu0htcCwIp6
DB_NAME=neondb
DB_SSL_MODE=require

# Kafka settings
KAFKA_BROKERS=localhost:9092
KAFKA_PAYMENTS_TOPIC=payments
KAFKA_SETTLEMENT_TOPIC=settlements
KAFKA_FRAUD_TOPIC=fraud
KAFKA_CONSUMER_GROUP=settlement-engine

# Settlement Engine settings
SETTLEMENT_CYCLE_INTERVAL_HOURS=24
DEFAULT_FEE_PERCENT=2.9
MINIMUM_SETTLEMENT_AMOUNT=100
```

## Database Connection Details

The service connects to a PostgreSQL database with the following connection string:

```
postgresql://neondb_owner:npg_QMu0htcCwIp6@ep-rapid-dew-a1roz96h-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require
```

## Running the Service

1. Ensure you have Go 1.20 or later installed
2. Make sure the `.env` file is set up correctly
3. Start the required infrastructure (Kafka, etc.) using Docker Compose:
   ```
   cd /path/to/fortexa
   docker-compose up -d
   ```
4. Run the Settlement Engine service:
   ```
   cd /path/to/fortexa/settlement-engine
   go run cmd/main.go
   ```

## Core Features

- **Payment Processing**: Marks captured payments as eligible for settlement
- **Settlement Cycle**: Runs at configurable intervals (default: 24 hours)
- **Fee Calculation**: Applies merchant-specific fee structures
- **Tax Processing**: Calculates applicable taxes on fees
- **Settlement Notifications**: Publishes settlement events to Kafka

## Settlement Process

1. The service consumes payment events from Kafka
2. When a payment is captured, it's marked as eligible for settlement
3. At scheduled intervals, the service creates settlement batches for eligible payments
4. Settlements are grouped by merchant and currency
5. The service calculates fees and taxes for each settlement
6. Settlement events are published back to Kafka for further processing

## Architecture

The Settlement Engine follows an event-driven architecture:

- **Kafka Consumer**: Reads payment events from the payments topic
- **Settlement Processor**: Contains the core business logic for settlement processing
- **Settlement Handler**: Manages settlement cycles and event processing
- **Kafka Producer**: Publishes settlement events to the settlements topic 
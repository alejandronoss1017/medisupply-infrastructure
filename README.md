# medisupply-infrastructure

This repository contains infrastructure and microservices for the MediSupply system, organized into two workshops (talleres) demonstrating different architectural patterns.

## Table of Contents
- [Taller 1: Synchronous Microservices Architecture](#taller-1-synchronous-microservices-architecture)
- [Taller 2: Event-Driven Microservices Architecture](#taller-2-event-driven-microservices-architecture)

---

## Taller 1: Synchronous Microservices Architecture

### Overview
Taller 1 demonstrates a **synchronous microservices architecture** where services communicate directly with each other via HTTP REST APIs.

### Technologies Used
- **Docker & Docker Compose**: Container orchestration
- **Node.js Microservices**: Business logic implementation
- **HTTP/REST**: Synchronous communication protocol
- **Health Checks**: Service availability monitoring

### Architecture Components
The system consists of 7 interconnected microservices:

1. **ruta-ms** (Port 3001): Route management service
   - Depends on: normativa-ms, vehiculo-ms
   
2. **centro-distribucion-ms** (Port 3002): Distribution center management
   
3. **lote-ms** (Port 3003): Batch/lot management
   - Depends on: normativa-ms, centro-distribucion-ms
   
4. **vehiculo-ms** (Port 3004): Vehicle management
   
5. **normativa-ms** (Port 3005): Regulatory compliance service
   - Depends on: alerta-ms
   
6. **alerta-ms** (Port 3006): Alert/notification service
   
7. **venta-ms** (Port 3007): Sales management service
   - Depends on: centro-distribucion-ms, ruta-ms

### How to Run

#### Prerequisites
- Docker and Docker Compose installed
- Ports 3001-3007 available

#### Steps

1. Navigate to the micro-dummies directory:
   ```bash
   cd micro-dummies
   ```

2. Start all services:
   ```bash
   docker-compose -f docker-compose.taller_1.yml up -d
   ```

3. Wait for all services to be healthy (check with):
   ```bash
   docker-compose -f docker-compose.taller_1.yml ps
   ```

4. To stop all services:
   ```bash
   docker-compose -f docker-compose.taller_1.yml down
   ```

### How to Test

#### Health Check
You can verify that all services are running by checking their health endpoints:
```bash
curl http://localhost:3001/health  # ruta-ms
curl http://localhost:3002/health  # centro-distribucion-ms
curl http://localhost:3003/health  # lote-ms
curl http://localhost:3004/health  # vehiculo-ms
curl http://localhost:3005/health  # normativa-ms
curl http://localhost:3006/health  # alerta-ms
curl http://localhost:3007/health  # venta-ms
```

#### Running the Integrated Tests
The docker-compose includes a tester service that exercises three main flows:

```bash
docker-compose -f docker-compose.taller_1.yml --profile test up tester
```

This will execute:
1. **Register Sale Flow**: POST to `/register-sale` endpoint
2. **Register Lot Flow**: POST to `/register-lot` endpoint
3. **Track Cold Chain Traceability Flow**: POST to `/track-cold-chain-traceability` endpoint

---

## Taller 2: Event-Driven Microservices Architecture

### Overview
Taller 2 demonstrates an **event-driven microservices architecture** using Kafka for asynchronous communication and DynamoDB for data persistence.

### Technologies Used
- **Docker & Docker Compose**: Container orchestration
- **Apache Kafka (KRaft mode)**: Event streaming platform for asynchronous messaging
- **DynamoDB Local**: NoSQL database for data persistence
- **Node.js Microservices**: Business logic implementation
- **Express.js**: Web framework for REST API
- **Kafka UI**: Web interface for monitoring Kafka topics and messages (Port 8080)
- **DynamoDB Admin**: Web interface for viewing DynamoDB tables and data (Port 8001)

### Architecture Components

Taller 2 is composed of **two docker-compose files**:

#### 1. Infrastructure (`docker-compose.infrastructure_taller_2.yml`)
Provides the foundational services:

- **Kafka** (Ports 9092, 29092, 9093): Message broker using KRaft consensus
- **DynamoDB Local** (Port 8000): Local DynamoDB instance for development
- **DynamoDB Admin** (Port 8001): Web UI for managing DynamoDB tables
- **Kafka UI** (Port 8080): Web UI for monitoring Kafka topics and events

#### 2. Applications (`docker-compose.apps_taller_2.yml`)
Contains the microservices:

- **suppliers-web** (Port 3001): REST API for supplier management
  - Publishes events to Kafka topic: `supplier-events`
  
- **suppliers-worker**: Consumes `supplier-events` (demonstration consumer)
  
- **purchase-plans-worker**: Consumes `supplier-events`
  - Stores data in DynamoDB table: `purchase-plans-db`
  
- **contracts-worker**: Consumes `supplier-events`
  - Stores data in DynamoDB table: `contracts-db`

### Event Flow
1. Client sends POST request to `suppliers-web` API
2. `suppliers-web` publishes event to Kafka topic `supplier-events`
3. Three workers consume the event in parallel:
   - `suppliers-worker`: Logs the event
   - `purchase-plans-worker`: Creates purchase plan in DynamoDB
   - `contracts-worker`: Creates contract in DynamoDB

### How to Run

#### Prerequisites
- Docker and Docker Compose installed
- Ports 3001, 8000, 8001, 8080, 9092 available

#### Steps

1. Navigate to the micro-dummies directory:
   ```bash
   cd micro-dummies
   ```

2. **Start the infrastructure first** (Kafka, DynamoDB, and UIs):
   ```bash
   docker-compose -f docker-compose.infrastructure_taller_2.yml up -d
   ```

3. Wait for infrastructure to be ready (especially Kafka):
   ```bash
   docker-compose -f docker-compose.infrastructure_taller_2.yml ps
   ```

4. **Start the applications**:
   ```bash
   docker-compose -f docker-compose.apps_taller_2.yml up -d
   ```

5. Verify all services are running:
   ```bash
   docker-compose -f docker-compose.infrastructure_taller_2.yml ps
   docker-compose -f docker-compose.apps_taller_2.yml ps
   ```

6. To stop all services:
   ```bash
   docker-compose -f docker-compose.apps_taller_2.yml down
   docker-compose -f docker-compose.infrastructure_taller_2.yml down
   ```

### How to Test

#### 1. Create a Supplier via REST API

Send a POST request to create a new supplier:

```bash
curl -X POST http://localhost:3001/api/suppliers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PROVEEDOR VACUNAS COVID"
  }'
```

Expected response: Success message with the created supplier information.

#### 2. Verify Event in Kafka UI

1. Open Kafka UI in your browser: **http://localhost:8080/**
2. Navigate to the `supplier-events` topic
3. Check the **Messages** tab
4. You should see the published event with the supplier data

#### 3. Verify Data in DynamoDB

##### Option A: Using DynamoDB Admin UI
1. Open DynamoDB Admin in your browser: **http://localhost:8001/**
2. You should see two tables:
   - `contracts-db`: Contains contract records created by contracts-worker
   - `purchase-plans-db`: Contains purchase plan records created by purchase-plans-worker
3. Click on each table to view the stored records

##### Option B: Using AWS CLI
If you have AWS CLI installed, you can query the tables directly:

```bash
# List tables
aws dynamodb list-tables --endpoint-url http://localhost:8000

# Scan contracts table
aws dynamodb scan --table-name contracts-db --endpoint-url http://localhost:8000

# Scan purchase-plans table
aws dynamodb scan --table-name purchase-plans-db --endpoint-url http://localhost:8000
```

#### 4. View Application Logs

To see the workers processing events in real-time:

```bash
# View all app logs
docker-compose -f docker-compose.apps_taller_2.yml logs -f

# View specific worker logs
docker-compose -f docker-compose.apps_taller_2.yml logs -f suppliers-worker
docker-compose -f docker-compose.apps_taller_2.yml logs -f purchase-plans-worker
docker-compose -f docker-compose.apps_taller_2.yml logs -f contracts-worker
```

### Troubleshooting

#### Services not connecting to Kafka
- Ensure infrastructure is fully started before launching apps
- Check Kafka health: `docker-compose -f docker-compose.infrastructure_taller_2.yml logs kafka`

#### DynamoDB tables not visible
- Verify DynamoDB Local is running: `docker ps | grep dynamodb-local`
- Check worker logs for table creation: `docker-compose -f docker-compose.apps_taller_2.yml logs contracts-worker`

#### Port conflicts
- Ensure no other services are using ports 3001, 8000, 8001, 8080, 9092
- Use `lsof -i :<port>` to check port usage

---

## License

This project is for educational purposes.

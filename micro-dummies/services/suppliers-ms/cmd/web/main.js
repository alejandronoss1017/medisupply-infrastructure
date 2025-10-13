// Main entry point for web server
const express = require('express');
const morgan = require('morgan');
const { SupplierServiceImpl } = require('../../internal/core/application/service');
const { KafkaProducer } = require('../../internal/adapter/kafka/producer');
const { SupplierHandler } = require('../../internal/adapter/http/handler');

const SERVICE_NAME = process.env.SERVICE_NAME || 'SUPPLIERS-MS';
const PORT = process.env.PORT || 3001;
const KAFKA_BROKERS = (process.env.KAFKA_BROKERS || 'localhost:9092').split(',');
const KAFKA_CLIENT_ID = process.env.KAFKA_CLIENT_ID || 'suppliers-ms';

async function main() {
    try {
        console.log(`Starting ${SERVICE_NAME}...`);
        console.log(`Kafka brokers: ${KAFKA_BROKERS.join(', ')}`);

        // Initialize Kafka producer
        const kafkaConfig = {
            clientId: KAFKA_CLIENT_ID,
            brokers: KAFKA_BROKERS
        };
        const kafkaProducer = new KafkaProducer(kafkaConfig);
        await kafkaProducer.connect();

        // Initialize service
        const supplierService = new SupplierServiceImpl(kafkaProducer);

        // Initialize HTTP handler
        const supplierHandler = new SupplierHandler(supplierService);

        // Setup Express app
        const app = express();
        app.use(express.json());
        app.use(morgan('dev'));

        // Register supplier routes
        app.use('/', supplierHandler.getRouter());

        // Start server
        const server = app.listen(PORT, () => {
            console.log(`${SERVICE_NAME} listening on port ${PORT}`);
        });

        // Graceful shutdown
        const gracefulShutdown = async () => {
            console.log('\nShutting down gracefully...');
            server.close(async () => {
                await kafkaProducer.close();
                console.log('Server closed');
                process.exit(0);
            });
        };

        process.on('SIGTERM', gracefulShutdown);
        process.on('SIGINT', gracefulShutdown);

    } catch (error) {
        console.error('Failed to start server:', error);
        process.exit(1);
    }
}

main();

// Main entry point for web server
const express = require('express');
const morgan = require('morgan');
const { PurchasePlanServiceImpl } = require('../../internal/core/application/service');
const { DynamoDBPurchasePlanRepository } = require('../../internal/adapter/dynamodb/repository');
const { PurchasePlanHandler } = require('../../internal/adapter/http/handler');

const SERVICE_NAME = process.env.SERVICE_NAME || 'PURCHASE-PLANS-MS';
const PORT = process.env.PORT || 3004;
const AWS_REGION = process.env.AWS_REGION || 'us-east-1';
const DYNAMODB_TABLE = process.env.DYNAMODB_TABLE || 'purchase-plans';

async function main() {
    try {
        console.log(`Starting ${SERVICE_NAME}...`);
        console.log(`DynamoDB region: ${AWS_REGION}`);
        console.log(`DynamoDB table: ${DYNAMODB_TABLE}`);

        // Initialize DynamoDB repository
        const repositoryConfig = {
            region: AWS_REGION,
            tableName: DYNAMODB_TABLE
        };
        const repository = new DynamoDBPurchasePlanRepository(repositoryConfig);

        // Initialize service
        const purchasePlanService = new PurchasePlanServiceImpl(repository);

        // Initialize HTTP handler
        const purchasePlanHandler = new PurchasePlanHandler(purchasePlanService);

        // Setup Express app
        const app = express();
        app.use(express.json());
        app.use(morgan('dev'));

        // Health check endpoint
        app.get('/purchase-plans/health', (_, res) => {
            res.json({ 
                status: 'ok', 
                service: SERVICE_NAME, 
                time: new Date().toISOString() 
            });
        });

        // Register purchase plan routes
        app.use('/', purchasePlanHandler.getRouter());

        // Start server
        const server = app.listen(PORT, () => {
            console.log(`${SERVICE_NAME} listening on port ${PORT}`);
        });

        // Graceful shutdown
        const gracefulShutdown = async () => {
            console.log('\nShutting down gracefully...');
            server.close(() => {
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

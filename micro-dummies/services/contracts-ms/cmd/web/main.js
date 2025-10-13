// Main entry point for web server
const express = require('express');
const morgan = require('morgan');
const { ContractServiceImpl } = require('../../internal/core/application/service');
const { DynamoDBRepository } = require('../../internal/adapter/dynamodb/repository');
const { ContractHandler } = require('../../internal/adapter/http/handler');

const SERVICE_NAME = process.env.SERVICE_NAME || 'CONTRACTS-MS';
const PORT = process.env.PORT || 3003;
const AWS_REGION = process.env.AWS_REGION || 'us-east-1';
const DYNAMODB_TABLE = process.env.DYNAMODB_TABLE || 'contracts';

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
        const repository = new DynamoDBRepository(repositoryConfig);

        // Initialize service
        const contractService = new ContractServiceImpl(repository);

        // Initialize HTTP handler
        const contractHandler = new ContractHandler(contractService);

        // Setup Express app
        const app = express();
        app.use(express.json());
        app.use(morgan('dev'));

        // Health check endpoint
        app.get('/contracts/health', (_, res) => {
            res.json({ 
                status: 'ok', 
                service: SERVICE_NAME, 
                time: new Date().toISOString() 
            });
        });

        // Register contract routes
        app.use('/', contractHandler.getRouter());

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

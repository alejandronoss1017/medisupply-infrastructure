// Main entry point for web server
const express = require('express');
const morgan = require('morgan');
const { PurchasePlanServiceImpl } = require('../../internal/core/application/service');
const { InMemoryPurchasePlanRepository } = require('../../internal/adapter/memory/repository');
const { PurchasePlanHandler } = require('../../internal/adapter/http/handler');

const SERVICE_NAME = process.env.SERVICE_NAME || 'PURCHASE-PLANS-MS';
const PORT = process.env.PORT || 3004;

async function main() {
    try {
        console.log(`Starting ${SERVICE_NAME}...`);

        // Initialize in-memory repository
        const repository = new InMemoryPurchasePlanRepository();

        // Initialize service
        const purchasePlanService = new PurchasePlanServiceImpl(repository);

        // Initialize HTTP handler
        const purchasePlanHandler = new PurchasePlanHandler(purchasePlanService);

        // Setup Express app
        const app = express();
        app.use(express.json());
        app.use(morgan('dev'));

        // Health check endpoint
        app.get('/health', (_, res) => {
            res.json({
                status: 'ok',
                service: process.env.SERVICE_NAME || 'PURCHASE-PLANS-MS',
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

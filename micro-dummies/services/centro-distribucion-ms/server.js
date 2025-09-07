const express = require('express');
const morgan = require('morgan');

const { logEvent } = require('./ddb-logger');
const DDB_TABLE = process.env.DDB_TABLE || 'distribution-center-db';

const app = express();
const SERVICE_NAME = process.env.SERVICE_NAME || 'CENTRO DE DISTRIBUCION MS';
const PORT = process.env.PORT || 3000;

app.use(express.json());
app.use(morgan('dev'));

app.get('/health', (_, res) => res.json({ status:'ok', service:SERVICE_NAME, time:new Date().toISOString() }));

app.all('/store-received-product', async (req, res) => {
    const result = { received: req.body || {} };
    await logEvent({ table: DDB_TABLE, service: SERVICE_NAME, endpoint: '/store-received-product', req, result });
    res.json({ service: SERVICE_NAME, endpoint: '/store-received-product', ...result, time: new Date().toISOString() });
});

app.all('/deduct-product-from-stock', async (req, res) => {
    const result = { deducted: req.body || {} };
    await logEvent({ table: DDB_TABLE, service: SERVICE_NAME, endpoint: '/deduct-product-from-stock', req, result });
    res.json({ service: SERVICE_NAME, endpoint: '/deduct-product-from-stock', ...result, time: new Date().toISOString() });
});

app.listen(PORT, () => console.log(`${SERVICE_NAME} listening on ${PORT}`));

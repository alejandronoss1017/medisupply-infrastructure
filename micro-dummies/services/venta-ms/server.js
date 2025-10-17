// server.js
/* eslint-disable no-console */
const express = require('express');
const morgan = require('morgan');
const axios = require('axios');
const { CloudEvent } = require('cloudevents');
const { Pool } = require('pg');

const { logEvent } = require('./ddb-logger');

// ======= Config =======
const DDB_TABLE = process.env.DDB_TABLE || 'sale-db';
const SERVICE_NAME = process.env.SERVICE_NAME || 'VENTA MS';
const PORT = +(process.env.PORT || 3000);

const URLs = {
    CENTRO_MS_URL: process.env.CENTRO_MS_URL || 'http://centro-distribucion-ms',
    RUTA_MS_URL: process.env.RUTA_MS_URL || 'http://ruta-ms',
};

// Postgres (primary write)
const PGHOST = process.env.PGHOST || 'pg-postgresql';
const PGPORT = +(process.env.PGPORT || 5432);
const PGUSER = process.env.PGUSER || 'app_user';
const PGPASSWORD = process.env.PGPASSWORD || 'password';
const PGDATABASE = process.env.PGDATABASE || 'salesdb';

// Read host (replica service)
const PGREADHOST = process.env.PGREADHOST || 'pg-postgresql-read';

// Pool options
const basePoolOpts = {
    port: PGPORT,
    user: PGUSER,
    password: PGPASSWORD,
    database: PGDATABASE,
    max: +(process.env.PG_POOL_MAX || 10),
    idleTimeoutMillis: +(process.env.PG_POOL_IDLE || 30000),
    connectionTimeoutMillis: +(process.env.PG_CONN_TIMEOUT || 5000),
};

// Primary (writes)
const pool = new Pool({ host: PGHOST, ...basePoolOpts });

// Replica (reads). If PGREADHOST not set, fallback to primary host.
const readPool = new Pool({ host: PGREADHOST || PGHOST, ...basePoolOpts });

// ======= Helpers =======
const FWD = [
    'x-request-id',
    'x-b3-traceid',
    'x-b3-spanid',
    'x-b3-parentspanid',
    'x-b3-sampled',
    'x-b3-flags',
    'x-ot-span-context',
    'traceparent',
    'tracestate',
];

const fwdHeaders = (req) => {
    const h = {};
    FWD.forEach((k) => { if (req.headers[k]) h[k] = req.headers[k]; });
    if (!h['x-request-id']) h['x-request-id'] = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
    return h;
};

async function callService(base, path, req, { method = 'post', data = {}, params = {} } = {}) {
    try {
        const resp = await axios({
            url: `${base}${path}`,
            method,
            data,
            params,
            headers: fwdHeaders(req),
            timeout: 3000,
            validateStatus: () => true,
        });
        return { ok: true, url: `${base}${path}`, status: resp.status, data: resp.data };
    } catch (e) {
        return { ok: false, url: `${base}${path}`, error: e.message };
    }
}

// Ensure table exists (idempotent) — do this only on PRIMARY
async function ensureSchema() {
    await pool.query(`
        CREATE TABLE IF NOT EXISTS sales (
                                             id TEXT PRIMARY KEY,
                                             sku TEXT NOT NULL,
                                             amount INT NOT NULL,
                                             created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            )
    `);
}

// Persist sale (PRIMARY)
async function saveSaleToDb({ saleId, sku, amount }) {
    await ensureSchema();
    await pool.query(
        `INSERT INTO sales (id, sku, amount) VALUES ($1, $2, $3)
            ON CONFLICT (id) DO NOTHING`,
        [saleId, sku, amount]
    );
}

// ======= App =======
const app = express();
app.use(express.json());
app.use(morgan('dev'));

// Basic health (app only)
app.get('/health', (_, res) => {
    res.json({ status: 'ok', service: SERVICE_NAME, time: new Date().toISOString() });
});

// Readiness/liveness (checks PRIMARY DB connection)
app.get('/ready', async (_, res) => {
    try {
        await pool.query('SELECT 1');
        res.json({ status: 'ready', pg: 'ok', service: SERVICE_NAME, time: new Date().toISOString() });
    } catch (e) {
        res.status(503).json({ status: 'not-ready', pg: e.message });
    }
});

// FIRST FLOW (parte 1): VENTA -> CENTRO + RUTA -> Persist PG + Log DDB
app.all('/register-sale', async (req, res) => {
    const sku = req.body?.sku || 'SKU-1';
    const amount = req.body?.amount ?? 1;
    const saleId = `S-${Date.now()}`;

    const chain = [];
    const payload = { sku, amount };

    chain.push(
        await callService(URLs.CENTRO_MS_URL, '/deduct-product-from-stock', req, {
            method: 'post',
            data: payload,
        })
    );
    chain.push(
        await callService(URLs.RUTA_MS_URL, '/plan-delivery-route', req, {
            method: 'post',
            data: { saleId, sku },
        })
    );

    // Guardar en Postgres (no falla el flujo si PG falla: se registra el error y se continúa)
    let pgStored = false;
    try {
        await saveSaleToDb({ saleId, sku, amount });
        pgStored = true;
    } catch (e) {
        console.error('PG insert error:', e);
    }

    const result = { chain, saleId, sku, amount, pgStored };
    try {
        await logEvent({ table: DDB_TABLE, service: SERVICE_NAME, endpoint: '/register-sale', req, result });
    } catch (e) {
        console.error('DDB logEvent error:', e);
    }

    res.json({ service: SERVICE_NAME, endpoint: '/register-sale', chain, saleId, pgStored, time: new Date().toISOString() });
});

// CloudEvent endpoint for Knative triggers
app.post('/events', async (req, res) => {
    try {
        const cloudEvent = CloudEvent.fromHTTP(req.headers, req.body);

        console.log('=== CloudEvent Received ===');
        console.log(`Event ID: ${cloudEvent.id}`);
        console.log(`Event Type: ${cloudEvent.type}`);
        console.log(`Event Source: ${cloudEvent.source}`);
        console.log(`Event Subject: ${cloudEvent.subject || 'N/A'}`);
        console.log(`Event Time: ${cloudEvent.time || 'N/A'}`);
        console.log(`Event Data Content Type: ${cloudEvent.datacontenttype || 'N/A'}`);
        if (cloudEvent.data) console.log('Event Data:', JSON.stringify(cloudEvent.data, null, 2));
        console.log('Event Extensions:');
        Object.entries(cloudEvent).forEach(([key, value]) => {
            if (!['id', 'type', 'source', 'subject', 'time', 'datacontenttype', 'data'].includes(key)) {
                console.log(`  ${key}: ${value}`);
            }
        });
        console.log('CloudEvent Headers:');
        Object.entries(req.headers).forEach(([key, value]) => console.log(`  ${key}: ${value}`));
        console.log('=== End CloudEvent ===');

        await logEvent({
            table: DDB_TABLE,
            service: SERVICE_NAME,
            endpoint: '/events',
            req,
            result: {
                cloudEvent: {
                    id: cloudEvent.id,
                    type: cloudEvent.type,
                    source: cloudEvent.source,
                    data: cloudEvent.data,
                },
            },
        });

        res.json({
            message: 'CloudEvent received and logged successfully',
            eventId: cloudEvent.id,
            eventType: cloudEvent.type,
            service: SERVICE_NAME,
            time: new Date().toISOString(),
        });
    } catch (error) {
        console.error('Failed to process CloudEvent:', error);
        res.status(400).json({
            error: 'Invalid CloudEvent format',
            message: error.message,
        });
    }
});

// Simple listing endpoint — READS FROM REPLICA (fallback to primary if needed)
app.get('/sales', async (req, res) => {
    try {
        const limit = Math.min(parseInt(req.query.limit || '50', 10), 200);
        const r = await readPool.query(
            'SELECT id, sku, amount, created_at FROM sales ORDER BY created_at DESC LIMIT $1',
            [limit]
        );
        res.json({ count: r.rowCount, items: r.rows, source: 'replica' });
    } catch (e) {
        console.error('Replica select error, falling back to primary:', e.message);
        try {
            const limit = Math.min(parseInt(req.query.limit || '50', 10), 200);
            const r = await pool.query(
                'SELECT id, sku, amount, created_at FROM sales ORDER BY created_at DESC LIMIT $1',
                [limit]
            );
            res.json({ count: r.rowCount, items: r.rows, source: 'primary' });
        } catch (e2) {
            console.error('Primary select error:', e2);
            res.status(500).json({ error: 'PG select failed', message: e2.message });
        }
    }
});

// Error middleware (last)
app.use((err, _req, res, _next) => {
    console.error('Unhandled error:', err);
    res.status(500).json({ error: 'internal_error', message: err?.message || 'unknown' });
});

// Graceful shutdown
function shutdown(signal) {
    console.log(`[${signal}] Shutting down...`);
    Promise.all([pool.end(), readPool.end()])
        .then(() => {
            console.log('PG pools closed');
            process.exit(0);
        })
        .catch((e) => {
            console.error('Error closing PG pools:', e);
            process.exit(1);
        });
}
['SIGINT', 'SIGTERM'].forEach((sig) => process.on(sig, () => shutdown(sig)));

app.listen(PORT, () => {
    console.log(`${SERVICE_NAME} listening on ${PORT}`);
    console.log(`Write DB → ${PGUSER}@${PGHOST}:${PGPORT}/${PGDATABASE}`);
    if (PGREADHOST) console.log(`Read DB → ${PGUSER}@${PGREADHOST}:${PGPORT}/${PGDATABASE}`);
});

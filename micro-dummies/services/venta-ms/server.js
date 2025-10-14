const express = require('express');
const morgan = require('morgan');
const axios = require('axios');
const { CloudEvent } = require('cloudevents');

const { logEvent } = require('./ddb-logger');
const DDB_TABLE = process.env.DDB_TABLE || 'sale-db';

const app = express();
const SERVICE_NAME = process.env.SERVICE_NAME || 'VENTA MS';
const PORT = process.env.PORT || 3000;

const URLs = {
    CENTRO_MS_URL: process.env.CENTRO_MS_URL || 'http://centro-distribucion-ms',
    RUTA_MS_URL: process.env.RUTA_MS_URL || 'http://ruta-ms'
};

const FWD = ['x-request-id','x-b3-traceid','x-b3-spanid','x-b3-parentspanid','x-b3-sampled','x-b3-flags','x-ot-span-context','traceparent','tracestate'];
const fwdHeaders = (req) => {
    const h = {}; FWD.forEach(k=>{ if(req.headers[k]) h[k]=req.headers[k]; });
    if(!h['x-request-id']) h['x-request-id'] = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
    return h;
};
async function callService(base, path, req, {method='post', data={}, params={}}={}) {
    try {
        const resp = await axios({ url: `${base}${path}`, method, data, params, headers: fwdHeaders(req), timeout: 3000, validateStatus: ()=>true });
        return { ok:true, url:`${base}${path}`, status:resp.status, data:resp.data };
    } catch (e) {
        return { ok:false, url:`${base}${path}`, error:e.message };
    }
}

app.use(express.json());
app.use(morgan('dev'));

app.get('/health', (_, res) => res.json({ status:'ok', service:SERVICE_NAME, time:new Date().toISOString() }));

// FIRST FLOW (parte 1): VENTA -> CENTRO + RUTA
app.all('/register-sale', async (req, res) => {
    const sku = req.body?.sku || 'SKU-1';
    const payload = { sku, amount: req.body?.amount ?? 1 };
    const chain = [];
    chain.push(await callService(URLs.CENTRO_MS_URL, '/deduct-product-from-stock', req, { method:'post', data: payload }));
    chain.push(await callService(URLs.RUTA_MS_URL, '/plan-delivery-route', req, { method:'post', data: { saleId: 'S-001', sku } }));
    const result = { chain };
    await logEvent({ table: DDB_TABLE, service: SERVICE_NAME, endpoint: '/register-sale', req, result });
    res.json({ service:SERVICE_NAME, endpoint:'/register-sale', chain, time:new Date().toISOString() });
});

// CloudEvent endpoint for Knative triggers
app.post('/events', async (req, res) => {
    try {
        // Parse the CloudEvent from the HTTP request
        const cloudEvent = CloudEvent.fromHTTP(req.headers, req.body);
        
        // Log the received CloudEvent details
        console.log('=== CloudEvent Received ===');
        console.log(`Event ID: ${cloudEvent.id}`);
        console.log(`Event Type: ${cloudEvent.type}`);
        console.log(`Event Source: ${cloudEvent.source}`);
        console.log(`Event Subject: ${cloudEvent.subject || 'N/A'}`);
        console.log(`Event Time: ${cloudEvent.time || 'N/A'}`);
        console.log(`Event Data Content Type: ${cloudEvent.datacontenttype || 'N/A'}`);
        
        // Log the event data
        if (cloudEvent.data) {
            console.log('Event Data:', JSON.stringify(cloudEvent.data, null, 2));
        }
        
        // Log all extensions/attributes
        console.log('Event Extensions:');
        Object.entries(cloudEvent).forEach(([key, value]) => {
            if (!['id', 'type', 'source', 'subject', 'time', 'datacontenttype', 'data'].includes(key)) {
                console.log(`  ${key}: ${value}`);
            }
        });
        
        // Log CloudEvent headers for debugging
        console.log('CloudEvent Headers:');
        Object.entries(req.headers).forEach(([key, value]) => {
            console.log(`  ${key}: ${value}`);
        });
        
        console.log('=== End CloudEvent ===');
        
        // Log the event to DynamoDB
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
                    data: cloudEvent.data
                }
            } 
        });
        
        // Respond with success
        res.json({
            message: 'CloudEvent received and logged successfully',
            eventId: cloudEvent.id,
            eventType: cloudEvent.type,
            service: SERVICE_NAME,
            time: new Date().toISOString()
        });
        
    } catch (error) {
        console.error('Failed to process CloudEvent:', error);
        res.status(400).json({ 
            error: 'Invalid CloudEvent format',
            message: error.message 
        });
    }
});

app.listen(PORT, () => console.log(`${SERVICE_NAME} listening on ${PORT}`));

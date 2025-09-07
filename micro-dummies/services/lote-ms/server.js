const express = require('express');
const morgan = require('morgan');
const axios = require('axios');

const { logEvent } = require('./ddb-logger');
const DDB_TABLE = process.env.DDB_TABLE || 'distribution-center-db';

const app = express();
const SERVICE_NAME = process.env.SERVICE_NAME || 'LOTE MS';
const PORT = process.env.PORT || 3000;

const URLs = {
    NORMATIVA_MS_URL: process.env.NORMATIVA_MS_URL || 'http://normativa-ms',
    CENTRO_MS_URL: process.env.CENTRO_MS_URL || 'http://centro-distribucion-ms'
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

// SECOND FLOW: LOTE -> NORMATIVA(validate-product-information) + CENTRO(store-received-product)
app.all('/register-lot', async (req, res) => {
    const lot = req.body?.lot || 'L-123';
    const chain = [];
    chain.push(await callService(URLs.NORMATIVA_MS_URL, '/validate-product-information', req, { method:'post', data:{ lot } }));
    chain.push(await callService(URLs.CENTRO_MS_URL, '/store-received-product', req, { method:'post', data:{ lot } }));
    const result = { lot: req.body?.lot || 'L-123', chain };
    await logEvent({ table: DDB_TABLE, service: SERVICE_NAME, endpoint: '/register-lot', req, result });
    res.json({ service: SERVICE_NAME, endpoint: '/register-lot', chain, time: new Date().toISOString() });
});

app.listen(PORT, () => console.log(`${SERVICE_NAME} listening on ${PORT}`));

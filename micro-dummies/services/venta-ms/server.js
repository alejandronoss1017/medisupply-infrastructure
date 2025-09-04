const express = require('express');
const morgan = require('morgan');
const axios = require('axios');

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
    res.json({ service:SERVICE_NAME, endpoint:'/register-sale', chain, time:new Date().toISOString() });
});

app.listen(PORT, () => console.log(`${SERVICE_NAME} listening on ${PORT}`));

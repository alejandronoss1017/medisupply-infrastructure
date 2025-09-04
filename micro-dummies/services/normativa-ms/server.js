const express = require('express');
const morgan = require('morgan');
const axios = require('axios');

const app = express();
const SERVICE_NAME = process.env.SERVICE_NAME || 'NORMATIVA MS';
const PORT = process.env.PORT || 3000;

const URLs = {
    ALERTA_MS_URL: process.env.ALERTA_MS_URL || 'http://alerta-ms'
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

// Usado por RUTA
app.all('/terms-of-delivery', (req, res) => {
    res.json({ service:SERVICE_NAME, endpoint:'/terms-of-delivery', terms:['CONDICION_A','CONDICION_B'], time:new Date().toISOString() });
});

// Usado por LOTE
app.all('/validate-product-information', (req, res) => {
    res.json({ service:SERVICE_NAME, endpoint:'/validate-product-information', valid:true, payload:req.body||{}, time:new Date().toISOString() });
});

// THIRD FLOW: NORMATIVA -> ALERTA
app.all('/track-cold-chain-traceability', async (req, res) => {
    const chain = [];
    chain.push(await callService(URLs.ALERTA_MS_URL, '/generate-alert', req, { method:'post', data:{ type:'COLD_CHAIN', sku:req.body?.sku || 'SKU-1' } }));
    res.json({ service:SERVICE_NAME, endpoint:'/track-cold-chain-traceability', chain, time:new Date().toISOString() });
});

app.listen(PORT, () => console.log(`${SERVICE_NAME} listening on ${PORT}`));

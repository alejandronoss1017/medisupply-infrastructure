const express = require('express');
const morgan = require('morgan');
const axios = require('axios');

const app = express();
const SERVICE_NAME = process.env.SERVICE_NAME || 'RUTA MS';
const PORT = process.env.PORT || 3000;

const URLs = {
  NORMATIVA_MS_URL: process.env.NORMATIVA_MS_URL || 'http://normativa-ms',
  VEHICULO_MS_URL: process.env.VEHICULO_MS_URL || 'http://vehiculo-ms'
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

// FIRST FLOW (parte 2): RUTA -> NORMATIVA(terms-of-delivery) + VEHICULO(assign-vehicle)
app.all('/plan-delivery-route', async (req, res) => {
  const chain = [];
  chain.push(await callService(URLs.NORMATIVA_MS_URL, '/terms-of-delivery', req, { method:'get', params:{ routeId: 'R-001' } }));
  chain.push(await callService(URLs.VEHICULO_MS_URL, '/assign-vehicle', req, { method:'post', data:{ routeId:'R-001' } }));
  res.json({ service:SERVICE_NAME, endpoint:'/plan-delivery-route', chain, time:new Date().toISOString() });
});

app.listen(PORT, () => console.log(`${SERVICE_NAME} listening on ${PORT}`));

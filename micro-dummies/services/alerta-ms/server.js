const express = require('express');
const morgan = require('morgan');

const app = express();
const SERVICE_NAME = process.env.SERVICE_NAME || 'ALERTA MS';
const PORT = process.env.PORT || 3000;

app.use(express.json());
app.use(morgan('dev'));

app.get('/health', (_, res) => res.json({ status:'ok', service:SERVICE_NAME, time:new Date().toISOString() }));

app.all('/generate-alert', (req, res) => {
    res.json({ service:SERVICE_NAME, endpoint:'/generate-alert', alert:req.body||{}, time:new Date().toISOString() });
});

app.listen(PORT, () => console.log(`${SERVICE_NAME} listening on ${PORT}`));

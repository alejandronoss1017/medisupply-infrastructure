const express = require('express');
const morgan = require('morgan');

const app = express();
const SERVICE_NAME = process.env.SERVICE_NAME || 'VEHICULO MS';
const PORT = process.env.PORT || 3000;

app.use(express.json());
app.use(morgan('dev'));

app.get('/health', (_, res) => res.json({ status:'ok', service:SERVICE_NAME, time:new Date().toISOString() }));

app.all('/assign-vehicle', (req, res) => {
    res.json({ service:SERVICE_NAME, endpoint:'/assign-vehicle', assignedTo: req.body?.routeId || 'R-UNK', time:new Date().toISOString() });
});

app.listen(PORT, () => console.log(`${SERVICE_NAME} listening on ${PORT}`));

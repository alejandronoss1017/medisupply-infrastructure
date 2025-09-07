// ddb-logger.js
const { DynamoDBClient } = require('@aws-sdk/client-dynamodb');
const { DynamoDBDocumentClient, PutCommand } = require('@aws-sdk/lib-dynamodb');
const { v4: uuidv4 } = require('uuid');

const TRACE_HEADERS = [
    'x-request-id','x-b3-traceid','x-b3-spanid','x-b3-parentspanid',
    'x-b3-sampled','x-b3-flags','x-ot-span-context','traceparent','tracestate'
];

const ddb = DynamoDBDocumentClient.from(new DynamoDBClient({
    region: process.env.AWS_REGION || 'us-east-1',
    credentials: (process.env.AWS_ACCESS_KEY_ID && process.env.AWS_SECRET_ACCESS_KEY) ? {
        accessKeyId: process.env.AWS_ACCESS_KEY_ID,
        secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
        sessionToken: process.env.AWS_SESSION_TOKEN
    } : undefined,
}));

function pick(obj, keys) {
    const out = {};
    for (const k of keys) if (obj[k] !== undefined) out[k] = obj[k];
    return out;
}

function truncate(obj, max = 4096) {
    try {
        const s = JSON.stringify(obj);
        if (s.length <= max) return obj;
        return { truncated: true, preview: s.slice(0, max) };
    } catch {
        return { note: 'non-serializable body' };
    }
}

function generateRandomInteger(digits) {
    const min = Math.pow(10, digits - 1);
    const max = Math.pow(10, digits) - 1;
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

async function logEvent({ table, service, endpoint, req, result, note }) {
    if (!table) return;
    const item = {
        id: generateRandomInteger(),
        ts: new Date().toISOString(),
        service,
        endpoint,
        method: req.method,
        trace: pick(req.headers || {}, TRACE_HEADERS),
        request: {
            path: req.path,
            query: req.query || {},
            body: truncate(req.body || {})
        },
        result: truncate(result || {}),
        pod: process.env.HOSTNAME || null,
        note: note || null
    };
    try {
        await ddb.send(new PutCommand({ TableName: table, Item: item }));
    } catch (e) {
        console.error(`[DDB] put failed (${table}):`, e.message);
    }
}

module.exports = { logEvent };

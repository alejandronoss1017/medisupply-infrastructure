// Main entry point for Kafka consumer worker
const { Kafka } = require('kafkajs'); // for admin API
const { KafkaConsumer } = require('../../internal/adapter/kafka/consumer');
const { PurchasePlanServiceImpl } = require('../../internal/core/application/service');
const { InMemoryPurchasePlanRepository } = require('../../internal/adapter/memory/repository');

const SERVICE_NAME = process.env.SERVICE_NAME || 'PURCHASE-PLANS-MS-WORKER';
const KAFKA_BROKERS = (process.env.KAFKA_BROKERS || 'localhost:9092').split(',');
const KAFKA_CLIENT_ID = process.env.KAFKA_CLIENT_ID || 'purchase-plans-ms-worker';
const KAFKA_GROUP_ID = process.env.KAFKA_GROUP_ID || 'purchase-plans-worker-group';
const KAFKA_TOPICS = (process.env.KAFKA_TOPICS || 'supplier-events').split(',').map(t => t.trim()).filter(Boolean);

// Removed DynamoDB configuration - using in-memory storage

const KAFKA_DEFAULT_PARTITIONS = parseInt(process.env.KAFKA_DEFAULT_PARTITIONS || '1', 10);
const KAFKA_DEFAULT_REPLICATION_FACTOR = parseInt(process.env.KAFKA_DEFAULT_REPLICATION_FACTOR || '1', 10);

let purchasePlanService;

async function ensureTopicsExist(brokers, topics, { numPartitions = 1, replicationFactor = 1 } = {}) {
    if (!topics.length) return;
    const kafka = new Kafka({ clientId: `${KAFKA_CLIENT_ID}-admin`, brokers });
    const admin = kafka.admin();
    await admin.connect();
    try {
        const existing = new Set(await admin.listTopics());
        const toCreate = topics.filter(t => !existing.has(t)).map(topic => ({ topic, numPartitions, replicationFactor }));
        if (toCreate.length) {
            console.log(`Creating missing topics: ${toCreate.map(t => t.topic).join(', ')}`);
            await admin.createTopics({ topics: toCreate, waitForLeaders: true });
            console.log('Topics created.');
        } else {
            console.log('All topics already exist.');
        }
    } finally {
        await admin.disconnect();
    }
}

async function handleMessage(topic, event) {
    try {
        console.log(`Processing message from topic ${topic}:`, event.eventType);
        const result = await purchasePlanService.processPurchasePlanEvent(event);
        // add event to event mesh
        //await addEventToEventMesh(event);
        if (result?.success) console.log('Purchase plan created successfully from event');
        else console.error('Failed to create purchase plan from event:', result?.error);
    } catch (error) {
        console.error('Error handling message:', error);
    }
}

async function addEventToEventMesh(event) {
    const http = require('http');
        //TODO: MOVE THIS TO WEB SERVICE
        // Adapt event into CloudEvent envelope
        // Example curl (for context):
        // curl -v "http://broker-ingress.knative-eventing.svc.cluster.local/order-system/order-broker" \
        //   -X POST \
        //   -H "Ce-Id: 12345" \
        //   -H "Ce-Specversion: 1.0" \
        //   -H "Ce-Type: order.created" \
        //   -H "Ce-Source: order-system/order-api" \
        //   -H "Content-Type: application/json" \
        //   -d '{ ... }'

        // Build CloudEvent headers and event body
        // You may want to customize these for your event type/domain
        const brokerHost = process.env.EVENTMESH_BROKER_HOST || 'broker-ingress.knative-eventing.svc.cluster.local';
        const brokerNamespace = process.env.EVENTMESH_BROKER_NAMESPACE || 'medisupply-broker';
        const brokerName = process.env.EVENTMESH_BROKER_NAME || 'medisupply-system';
        const ceType = event.eventType || 'updated.med';
        const ceId = event.id || String(Date.now());
        const ceSource = process.env.EVENTMESH_CE_SOURCE || 'purchase-plans-ms/worker';

        const bodyStr = JSON.stringify(event);
        const options = {
            hostname: brokerHost,
            port: 80,
            path: `/${brokerNamespace}/${brokerName}`,
            method: 'POST',
            headers: {
                'Ce-Id': ceId,
                'Ce-Specversion': '1.0',
                'Ce-Type': ceType,
                'Ce-Source': ceSource,
                'Content-Type': 'application/json',
                'Content-Length': Buffer.byteLength(bodyStr)
            }
        };

        await new Promise((resolve, reject) => {
            const req = http.request(options, (res) => {
                let data = '';
                res.on('data', chunk => data += chunk);
                res.on('end', () => {
                    if (res.statusCode >= 200 && res.statusCode < 300) {
                        console.log('Event pushed to event mesh:', res.statusCode);
                        resolve();
                    } else {
                        console.error('Event mesh responded with:', res.statusCode, data);
                        reject(new Error(`Event mesh HTTP error: ${res.statusCode}`));
                    }
                });
            });
            req.on('error', (err) => {
                console.error('Failed to push to event mesh:', err);
                reject(err);
            });
            req.write(bodyStr);
        req.end();
    });
    console.log('Adding event to event mesh:', event);
}

async function main() {
    try {
        console.log(`Starting ${SERVICE_NAME}...`);
        console.log(`Kafka brokers: ${KAFKA_BROKERS.join(', ')}`);
        console.log(`Subscribing to topics: ${KAFKA_TOPICS.join(', ')}`);

        // In-memory repository
        const repository = new InMemoryPurchasePlanRepository();
        purchasePlanService = new PurchasePlanServiceImpl(repository);

        // Ensure topics exist BEFORE subscribe
        await ensureTopicsExist(KAFKA_BROKERS, KAFKA_TOPICS, {
            numPartitions: KAFKA_DEFAULT_PARTITIONS,
            replicationFactor: KAFKA_DEFAULT_REPLICATION_FACTOR
        });

        // Kafka consumer
        const kafkaConsumer = new KafkaConsumer(
            { clientId: `${KAFKA_CLIENT_ID}-consumer`, brokers: KAFKA_BROKERS, groupId: KAFKA_GROUP_ID },
            handleMessage
        );

        await kafkaConsumer.subscribe(KAFKA_TOPICS);
        console.log('Starting to consume messages...');
        await kafkaConsumer.consume();

        const gracefulShutdown = async () => {
            console.log('\nShutting down gracefully...');
            await kafkaConsumer.close();
            console.log('Worker closed');
            process.exit(0);
        };
        process.on('SIGTERM', gracefulShutdown);
        process.on('SIGINT', gracefulShutdown);
    } catch (error) {
        console.error('Failed to start worker:', error);
        process.exit(1);
    }
}

main();

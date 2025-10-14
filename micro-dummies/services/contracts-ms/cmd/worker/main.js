// Main entry point for Kafka consumer worker
const { Kafka } = require('kafkajs');
const { KafkaConsumer } = require('../../internal/adapter/kafka/consumer');

const SERVICE_NAME = process.env.SERVICE_NAME || 'CONTRACTS-MS-WORKER';
const KAFKA_BROKERS = (process.env.KAFKA_BROKERS || 'localhost:9092').split(',');
const KAFKA_CLIENT_ID = process.env.KAFKA_CLIENT_ID || 'contracts-ms-worker';
const KAFKA_GROUP_ID = process.env.KAFKA_GROUP_ID || 'contracts-worker-group';
const KAFKA_TOPICS = (process.env.KAFKA_TOPICS || 'supplier-events').split(',').map(t => t.trim()).filter(Boolean);


const KAFKA_DEFAULT_PARTITIONS = parseInt(process.env.KAFKA_DEFAULT_PARTITIONS || '1', 10);
const KAFKA_DEFAULT_REPLICATION_FACTOR = parseInt(process.env.KAFKA_DEFAULT_REPLICATION_FACTOR || '1', 10);

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

async function sendCloudEvent(eventData, eventType, eventId = crypto.randomUUID()) {
    const response = await fetch(
        "http://broker-ingress.knative-eventing.svc.cluster.local/order-system/order-broker",
        {
            method: "POST",
            headers: {
                "Ce-Id": eventId,
                "Ce-Specversion": "1.0",
                "Ce-Type": eventType,
                "Ce-Source": "order-system/order-api",
                "Content-Type": "application/json"
            },
            body: JSON.stringify(eventData)
        }
    );

    if (!response.ok) {
        throw new Error(`Failed to send event: ${response.status} ${response.statusText}`);
    }

    return response;
}

async function handleMessage(topic, event) {

    console.log(`Processing message from topic ${topic}:`, event.eventType);

    try{
    const response =  await sendCloudEvent(event, event.eventType, event.eventId);

    console.log('Event mesh response status:', response.status);
    console.log('Event sent to event mesh successfully');

    }catch (e) {
        console.error('Failed to send event to event mesh:', e);
    }
}

async function main() {
    try {
        console.log(`Starting ${SERVICE_NAME}...`);
        console.log(`Kafka brokers: ${KAFKA_BROKERS.join(', ')}`);
        console.log(`Subscribing to topics: ${KAFKA_TOPICS.join(', ')}`);

        await ensureTopicsExist(KAFKA_BROKERS, KAFKA_TOPICS, {
            numPartitions: KAFKA_DEFAULT_PARTITIONS,
            replicationFactor: KAFKA_DEFAULT_REPLICATION_FACTOR
        });

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

// Main entry point for Kafka consumer worker
const { Kafka } = require('kafkajs'); // <-- for admin API to create topics
const { KafkaConsumer } = require('../../internal/adapter/kafka/consumer');
const { SupplierServiceImpl } = require('../../internal/core/application/service');
const { KafkaProducer } = require('../../internal/adapter/kafka/producer');

const SERVICE_NAME = process.env.SERVICE_NAME || 'SUPPLIERS-MS-WORKER';
const KAFKA_BROKERS = (process.env.KAFKA_BROKERS || 'localhost:9092').split(',');
const KAFKA_CLIENT_ID = process.env.KAFKA_CLIENT_ID || 'suppliers-ms-worker';
const KAFKA_GROUP_ID = process.env.KAFKA_GROUP_ID || 'suppliers-worker-group';
const KAFKA_TOPICS = (process.env.KAFKA_TOPICS || 'supplier-requests')
    .split(',')
    .map(t => t.trim())
    .filter(Boolean);

// optional tuning via env
const KAFKA_DEFAULT_PARTITIONS = parseInt(process.env.KAFKA_DEFAULT_PARTITIONS || '1', 10);
const KAFKA_DEFAULT_REPLICATION_FACTOR = parseInt(process.env.KAFKA_DEFAULT_REPLICATION_FACTOR || '1', 10);

async function ensureTopicsExist(brokers, topics, { numPartitions = 1, replicationFactor = 1 } = {}) {
    if (!topics.length) return;

    const kafka = new Kafka({ clientId: `${KAFKA_CLIENT_ID}-admin`, brokers });
    const admin = kafka.admin();
    await admin.connect();

    try {
        const existing = new Set(await admin.listTopics());
        const toCreate = topics
            .filter(t => !existing.has(t))
            .map(topic => ({ topic, numPartitions, replicationFactor }));

        if (toCreate.length > 0) {
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
        console.log(`Processing message from topic ${topic}:`, event);
        // Add custom message processing logic here
    } catch (error) {
        console.error('Error handling message:', error);
    }
}

async function main() {
    try {
        console.log(`Starting ${SERVICE_NAME}...`);
        console.log(`Kafka brokers: ${KAFKA_BROKERS.join(', ')}`);
        console.log(`Subscribing to topics: ${KAFKA_TOPICS.join(', ')}`);

        // Initialize Kafka producer (for republishing events)
        const kafkaProducerConfig = {
            clientId: KAFKA_CLIENT_ID,
            brokers: KAFKA_BROKERS
        };
        const kafkaProducer = new KafkaProducer(kafkaProducerConfig);
        await kafkaProducer.connect();

        // Initialize service
        const supplierService = new SupplierServiceImpl(kafkaProducer);

        // Ensure topics exist BEFORE creating consumer subscription
        await ensureTopicsExist(KAFKA_BROKERS, KAFKA_TOPICS, {
            numPartitions: KAFKA_DEFAULT_PARTITIONS,
            replicationFactor: KAFKA_DEFAULT_REPLICATION_FACTOR
        });

        // Initialize Kafka consumer
        const kafkaConsumerConfig = {
            clientId: KAFKA_CLIENT_ID + '-consumer',
            brokers: KAFKA_BROKERS,
            groupId: KAFKA_GROUP_ID
        };
        const kafkaConsumer = new KafkaConsumer(kafkaConsumerConfig, handleMessage);

        // Subscribe to topics
        await kafkaConsumer.subscribe(KAFKA_TOPICS);

        // Start consuming
        console.log('Starting to consume messages...');
        await kafkaConsumer.consume();

        // Graceful shutdown
        const gracefulShutdown = async () => {
            console.log('\nShutting down gracefully...');
            await kafkaConsumer.close();
            await kafkaProducer.close();
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

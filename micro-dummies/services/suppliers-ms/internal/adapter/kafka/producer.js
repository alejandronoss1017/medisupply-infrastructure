// Kafka producer adapter
const { Kafka } = require('kafkajs');
const { EventPublisher } = require('../../core/port/driven/publisher');

class KafkaProducer extends EventPublisher {
    constructor(config) {
        super();
        this.kafka = new Kafka({
            clientId: config.clientId || 'suppliers-ms',
            brokers: config.brokers || ['localhost:9092']
        });
        this.producer = this.kafka.producer();
        this.connected = false;
    }

    async connect() {
        if (!this.connected) {
            await this.producer.connect();
            this.connected = true;
            console.log('Kafka producer connected');
        }
    }

    async publish(topic, event) {
        try {
            if (!this.connected) {
                await this.connect();
            }

            const message = {
                key: event.eventId || Date.now().toString(),
                value: JSON.stringify(event),
                headers: {
                    'event-type': event.eventType || 'UNKNOWN',
                    'timestamp': event.timestamp || new Date().toISOString()
                }
            };

            await this.producer.send({
                topic,
                messages: [message]
            });

            console.log(`Event published to topic ${topic}:`, event.eventType);
            return true;
        } catch (error) {
            console.error('Error publishing event to Kafka:', error);
            throw error;
        }
    }

    async close() {
        if (this.connected) {
            await this.producer.disconnect();
            this.connected = false;
            console.log('Kafka producer disconnected');
        }
    }
}

module.exports = { KafkaProducer };

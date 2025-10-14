// Kafka consumer adapter
const { Kafka } = require('kafkajs');

class KafkaConsumer {
    constructor(config, messageHandler) {
        this.kafka = new Kafka({
            clientId: config.clientId || 'suppliers-ms-consumer',
            brokers: config.brokers || ['localhost:9092']
        });
        this.consumer = this.kafka.consumer({ 
            groupId: config.groupId || 'suppliers-group' 
        });
        this.messageHandler = messageHandler;
        this.connected = false;
    }

    async connect() {
        if (!this.connected) {
            await this.consumer.connect();
            this.connected = true;
            console.log('Kafka consumer connected');
        }
    }

    async subscribe(topics) {
        try {
            if (!this.connected) {
                await this.connect();
            }

            for (const topic of topics) {
                await this.consumer.subscribe({ topic, fromBeginning: false });
                console.log(`Subscribed to topic: ${topic}`);
            }
        } catch (error) {
            console.error('Error subscribing to topics:', error);
            throw error;
        }
    }

    async consume() {
        try {
            await this.consumer.run({
                eachMessage: async ({ topic, partition, message }) => {
                    try {
                        const value = message.value.toString();
                        const event = JSON.parse(value);
                        
                        console.log(`Received message from ${topic}:`, {
                            partition,
                            offset: message.offset,
                            eventType: event.eventType
                        });

                        if (this.messageHandler) {
                            await this.messageHandler(topic, event);
                        }
                    } catch (error) {
                        console.error('Error processing message:', error);
                    }
                }
            });
            // Keep the promise pending indefinitely to prevent process exit
            await new Promise(() => {});
        } catch (error) {
            console.error('Error in consumer run:', error);
            throw error;
        }
    }

    async close() {
        if (this.connected) {
            await this.consumer.disconnect();
            this.connected = false;
            console.log('Kafka consumer disconnected');
        }
    }
}

module.exports = { KafkaConsumer };

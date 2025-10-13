// Publisher port - interface for publishing events to message broker
class EventPublisher {
    async publish(topic, event) {
        throw new Error('publish method must be implemented');
    }

    async close() {
        throw new Error('close method must be implemented');
    }
}

module.exports = { EventPublisher };

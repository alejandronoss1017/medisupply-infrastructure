// Domain events for Supplier
class SupplierEvent {
    constructor(eventType, data) {
        this.eventType = eventType;
        this.data = data;
        this.timestamp = new Date().toISOString();
        this.eventId = this.generateEventId();
    }

    generateEventId() {
        return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    }

    toJSON() {
        return {
            eventId: this.eventId,
            eventType: this.eventType,
            data: this.data,
            timestamp: this.timestamp
        };
    }
}

class SupplierCreatedEvent extends SupplierEvent {
    constructor(supplier) {
        super('SUPPLIER_CREATED', supplier);
    }
}

class SupplierUpdatedEvent extends SupplierEvent {
    constructor(supplier) {
        super('SUPPLIER_UPDATED', supplier);
    }
}

module.exports = {
    SupplierEvent,
    SupplierCreatedEvent,
    SupplierUpdatedEvent
};

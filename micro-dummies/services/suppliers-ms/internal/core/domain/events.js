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
        super('supplier.created', supplier);
    }
}

class SupplierUpdatedEvent extends SupplierEvent {
    constructor(supplier) {
        super('supplier.updated', supplier);
    }
}

class MedicineEvent {
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

class MedicineCreatedEvent extends MedicineEvent {
    constructor(medicine) {
        super('medicine.updated', medicine);
    }
}

class MedicineUpdatedEvent extends MedicineEvent {
    constructor(medicine) {
        super('medicine.updated', medicine);
    }
}

module.exports = {
    SupplierEvent,
    SupplierCreatedEvent,
    SupplierUpdatedEvent,
    MedicineEvent,
    MedicineCreatedEvent,
    MedicineUpdatedEvent
};

// Application service - implements business logic
const { Supplier } = require('../domain/supplier');
const { SupplierCreatedEvent, SupplierUpdatedEvent } = require('../domain/events');
const { randomUUID } = require('node:crypto');

class SupplierServiceImpl {
    constructor(eventPublisher) {
        this.eventPublisher = eventPublisher;
    }

    async createSupplier(supplierData) {
        try {
            const supplier = new Supplier({
                id: randomUUID(),
                ...supplierData
            });

            supplier.validate();

            const event = new SupplierCreatedEvent(supplier.toJSON());
            await this.eventPublisher.publish('supplier-events', event.toJSON());

            return { success: true, supplier: supplier.toJSON(), eventId: event.eventId };
        } catch (error) {
            console.error('Error creating supplier:', error);
            throw error;
        }
    }

    async updateSupplier(id, supplierData) {
        try {
            const supplier = new Supplier({ id, ...supplierData });
            supplier.validate();

            const event = new SupplierUpdatedEvent(supplier.toJSON());
            await this.eventPublisher.publish('supplier-events', event.toJSON());

            return { success: true, supplier: supplier.toJSON(), eventId: event.eventId };
        } catch (error) {
            console.error('Error updating supplier:', error);
            throw error;
        }
    }

    async processSupplierRequest(requestData) {
        try {
            const supplier = new Supplier({
                id: requestData.id || randomUUID(),
                ...requestData
            });

            supplier.validate();

            const event = new SupplierCreatedEvent(supplier.toJSON());
            await this.eventPublisher.publish('supplier-events', event.toJSON());

            return {
                success: true,
                message: 'Supplier request processed successfully',
                supplier: supplier.toJSON(),
                eventId: event.eventId
            };
        } catch (error) {
            console.error('Error processing supplier request:', error);
            throw error;
        }
    }
}

module.exports = { SupplierServiceImpl };

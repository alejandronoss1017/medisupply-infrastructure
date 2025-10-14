// Application service - implements business logic
const { Supplier } = require('../domain/supplier');
const { Medicine } = require('../domain/medicine');
const { SupplierCreatedEvent, SupplierUpdatedEvent, MedicineCreatedEvent, MedicineUpdatedEvent } = require('../domain/events');
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
            await this.eventPublisher.publish('supplier.events', event.toJSON());

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
            await this.eventPublisher.publish('supplier.events', event.toJSON());

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

class MedicineServiceImpl {
    constructor(eventPublisher) {
        this.eventPublisher = eventPublisher;
    }

    async createMedicine(medicineData) {
        try {
            const medicine = new Medicine({
                id: randomUUID(),
                ...medicineData
            });

            medicine.validate();

            const event = new MedicineCreatedEvent(medicine.toJSON());
            await this.eventPublisher.publish('medicine.events', event.toJSON());

            return { success: true, medicine: medicine.toJSON(), eventId: event.eventId };
        } catch (error) {
            console.error('Error creating medicine:', error);
            throw error;
        }
    }

    async updateMedicine(id, medicineData) {
        try {
            const medicine = new Medicine({ id, ...medicineData });
            medicine.validate();

            const event = new MedicineUpdatedEvent(medicine.toJSON());
            await this.eventPublisher.publish('medicine.events', event.toJSON());

            return { success: true, medicine: medicine.toJSON(), eventId: event.eventId };
        } catch (error) {
            console.error('Error updating medicine:', error);
            throw error;
        }
    }

}


module.exports = { SupplierServiceImpl, MedicineServiceImpl };

// Application service - implements business logic
const { PurchasePlan } = require('../domain/purchasePlan');
const { randomUUID } = require('node:crypto');

class PurchasePlanServiceImpl {
    constructor(repository) {
        this.repository = repository;
    }

    async createPurchasePlan(planData) {
        try {
            const purchasePlan = new PurchasePlan({
                id: randomUUID(),
                ...planData
            });

            purchasePlan.validate();

            await this.repository.save(purchasePlan.toJSON());

            return {
                success: true,
                purchasePlan: purchasePlan.toJSON()
            };
        } catch (error) {
            console.error('Error creating purchase plan:', error);
            throw error;
        }
    }

    async processPurchasePlanEvent(event) {
        try {
            console.log('Processing purchase plan event:', event.eventType);

            // Extract supplier data from event
            const eventData = event?.data || {};

            // Create purchase plan from supplier event
            const purchasePlan = new PurchasePlan({
                id: randomUUID(),
                supplierId: eventData.id || randomUUID(),
                items: Array.isArray(eventData.items) ? eventData.items : [],
                totalAmount: typeof eventData.totalAmount === 'number' ? eventData.totalAmount : 0,
                status: 'created_from_event',
                metadata: {
                    sourceEvent: event.eventType,
                    sourceEventId: event.eventId,
                    sourceTimestamp: event.timestamp,
                    supplierData: eventData
                }
            });

            // Persist
            await this.repository.save(purchasePlan.toJSON());
            console.log('Purchase plan saved from event:', purchasePlan.id);

            return {
                success: true,
                purchasePlan: purchasePlan.toJSON()
            };
        } catch (error) {
            console.error('Error processing purchase plan event:', error);
            // Don't throw; keep consumer alive
            return {
                success: false,
                error: error.message
            };
        }
    }

    async getPurchasePlan(id) {
        try {
            return await this.repository.findById(id);
        } catch (error) {
            console.error('Error getting purchase plan:', error);
            throw error;
        }
    }

    async getAllPurchasePlans() {
        try {
            return await this.repository.findAll();
        } catch (error) {
            console.error('Error getting all purchase plans:', error);
            throw error;
        }
    }
}

module.exports = { PurchasePlanServiceImpl };

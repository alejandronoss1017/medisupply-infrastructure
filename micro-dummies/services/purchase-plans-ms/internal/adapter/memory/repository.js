// In-memory repository adapter
const { PurchasePlanRepository } = require('../../core/port/driven/repository');

class InMemoryPurchasePlanRepository extends PurchasePlanRepository {
    constructor() {
        super();
        this.data = new Map();
        console.log('In-memory repository initialized');
    }

    async save(purchasePlan) {
        try {
            const item = {
                ...purchasePlan,
                ts: new Date().toISOString()
            };

            this.data.set(purchasePlan.id, item);
            console.log(`Purchase plan saved to memory: ${purchasePlan.id}`);
            return item;
        } catch (error) {
            console.error('Error saving to memory:', error);
            throw error;
        }
    }

    async findById(id) {
        try {
            return this.data.get(id) || null;
        } catch (error) {
            console.error('Error finding by ID in memory:', error);
            throw error;
        }
    }

    async findAll() {
        try {
            return Array.from(this.data.values());
        } catch (error) {
            console.error('Error getting all items from memory:', error);
            throw error;
        }
    }

    async update(id, purchasePlan) {
        try {
            const item = {
                ...purchasePlan,
                id,
                updatedAt: new Date().toISOString()
            };

            this.data.set(id, item);
            return item;
        } catch (error) {
            console.error('Error updating in memory:', error);
            throw error;
        }
    }

    async delete(id) {
        try {
            const deleted = this.data.delete(id);
            console.log(`Purchase plan deleted from memory: ${id}`);
            return deleted;
        } catch (error) {
            console.error('Error deleting from memory:', error);
            throw error;
        }
    }
}

module.exports = { InMemoryPurchasePlanRepository };

// Service port - interface for purchase plan operations
class PurchasePlanService {
    async createPurchasePlan(planData) {
        throw new Error('createPurchasePlan method must be implemented');
    }

    async processPurchasePlanEvent(event) {
        throw new Error('processPurchasePlanEvent method must be implemented');
    }

    async getPurchasePlan(id) {
        throw new Error('getPurchasePlan method must be implemented');
    }

    async getAllPurchasePlans() {
        throw new Error('getAllPurchasePlans method must be implemented');
    }
}

module.exports = { PurchasePlanService };

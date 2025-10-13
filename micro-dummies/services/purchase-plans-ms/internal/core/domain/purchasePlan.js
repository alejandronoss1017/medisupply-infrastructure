// Domain model for Purchase Plan
class PurchasePlan {
    constructor({ id, supplierId, items, totalAmount, status, metadata = {} }) {
        this.id = id;
        this.supplierId = supplierId;
        this.items = items || [];
        this.totalAmount = totalAmount || 0;
        this.status = status || 'pending';
        this.metadata = metadata;
        this.createdAt = new Date().toISOString();
    }

    validate() {
        if (!this.supplierId || this.supplierId.trim() === '') {
            throw new Error('Supplier ID is required');
        }
        if (!this.items || this.items.length === 0) {
            throw new Error('At least one item is required');
        }
        return true;
    }

    toJSON() {
        return {
            id: this.id,
            supplierId: this.supplierId,
            items: this.items,
            totalAmount: this.totalAmount,
            status: this.status,
            metadata: this.metadata,
            createdAt: this.createdAt
        };
    }
}

module.exports = { PurchasePlan };

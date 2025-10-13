// Domain model for Contract
class Contract {
    constructor({ id, supplierId, terms, value, startDate, endDate, status, metadata = {} }) {
        this.id = id;
        this.supplierId = supplierId;
        this.terms = terms || '';
        this.value = value || 0;
        this.startDate = startDate;
        this.endDate = endDate;
        this.status = status || 'draft';
        this.metadata = metadata;
        this.createdAt = new Date().toISOString();
    }

    validate() {
        if (!this.supplierId || this.supplierId.trim() === '') {
            throw new Error('Supplier ID is required');
        }
        if (!this.terms || this.terms.trim() === '') {
            throw new Error('Contract terms are required');
        }
        if (this.value <= 0) {
            throw new Error('Contract value must be greater than zero');
        }
        return true;
    }

    toJSON() {
        return {
            id: this.id,
            supplierId: this.supplierId,
            terms: this.terms,
            value: this.value,
            startDate: this.startDate,
            endDate: this.endDate,
            status: this.status,
            metadata: this.metadata,
            createdAt: this.createdAt
        };
    }
}

module.exports = { Contract };

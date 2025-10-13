// Domain model for Supplier
class Supplier {
    constructor({ id, name, contactInfo, address, status, metadata = {} }) {
        this.id = id;
        this.name = name;
        this.contactInfo = contactInfo;
        this.address = address;
        this.status = status || 'active';
        this.metadata = metadata;
        this.createdAt = new Date().toISOString();
    }

    validate() {
        if (!this.name || this.name.trim() === '') {
            throw new Error('Supplier name is required');
        }
        return true;
    }

    toJSON() {
        return {
            id: this.id,
            name: this.name,
            contactInfo: this.contactInfo,
            address: this.address,
            status: this.status,
            metadata: this.metadata,
            createdAt: this.createdAt
        };
    }
}

module.exports = { Supplier };

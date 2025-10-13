// Service port - interface for supplier operations
class SupplierService {
    async createSupplier(supplierData) {
        throw new Error('createSupplier method must be implemented');
    }

    async updateSupplier(id, supplierData) {
        throw new Error('updateSupplier method must be implemented');
    }

    async processSupplierRequest(requestData) {
        throw new Error('processSupplierRequest method must be implemented');
    }
}

module.exports = { SupplierService };

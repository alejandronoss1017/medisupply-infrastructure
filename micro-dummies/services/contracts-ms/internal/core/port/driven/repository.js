// Repository port - interface for data persistence
class ContractRepository {
    async save(contract) {
        throw new Error('save method must be implemented');
    }

    async findById(id) {
        throw new Error('findById method must be implemented');
    }

    async findAll() {
        throw new Error('findAll method must be implemented');
    }

    async update(id, contract) {
        throw new Error('update method must be implemented');
    }

    async delete(id) {
        throw new Error('delete method must be implemented');
    }
}

module.exports = { ContractRepository };

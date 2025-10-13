// Service port - interface for contract operations
class ContractService {
    async createContract(contractData) {
        throw new Error('createContract method must be implemented');
    }

    async processContractEvent(event) {
        throw new Error('processContractEvent method must be implemented');
    }

    async getContract(id) {
        throw new Error('getContract method must be implemented');
    }

    async getAllContracts() {
        throw new Error('getAllContracts method must be implemented');
    }
}

module.exports = { ContractService };

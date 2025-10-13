// Application service - implements business logic
const { Contract } = require('../domain/contract');
const { randomUUID } = require('node:crypto');

class ContractServiceImpl {
    constructor(repository) {
        this.repository = repository;
    }

    async createContract(contractData) {
        try {
            const contract = new Contract({
                id: randomUUID(),
                ...contractData
            });

            contract.validate();

            await this.repository.save(contract.toJSON());

            return {
                success: true,
                contract: contract.toJSON()
            };
        } catch (error) {
            console.error('Error creating contract:', error);
            throw error;
        }
    }

    async processContractEvent(event) {
        try {
            console.log('Processing contract event:', event.eventType);

            // Extract supplier data from event
            const eventData = event?.data || {};

            // Create contract from supplier event
            const contract = new Contract({
                id: randomUUID(),
                supplierId: eventData.id || randomUUID(),
                terms: eventData.terms || `Auto-generated terms for supplier ${eventData.name || 'Unknown'}`,
                value:
                    typeof eventData.value === 'number'
                        ? eventData.value
                        : typeof eventData.totalAmount === 'number'
                            ? eventData.totalAmount
                            : 10000,
                startDate: new Date().toISOString(),
                endDate: this.calculateEndDate(12), // default 12 months
                status: 'created_from_event',
                metadata: {
                    sourceEvent: event.eventType,
                    sourceEventId: event.eventId,
                    sourceTimestamp: event.timestamp,
                    supplierData: eventData
                }
            });

            // Persist
            await this.repository.save(contract.toJSON());
            console.log('Contract saved from event:', contract.id);

            return {
                success: true,
                contract: contract.toJSON()
            };
        } catch (error) {
            console.error('Error processing contract event:', error);
            // Don't throw; keep consumer alive
            return {
                success: false,
                error: error.message
            };
        }
    }

    calculateEndDate(months) {
        const date = new Date();
        date.setMonth(date.getMonth() + months);
        return date.toISOString();
    }

    async getContract(id) {
        try {
            return await this.repository.findById(id);
        } catch (error) {
            console.error('Error getting contract:', error);
            throw error;
        }
    }

    async getAllContracts() {
        try {
            return await this.repository.findAll();
        } catch (error) {
            console.error('Error getting all contracts:', error);
            throw error;
        }
    }
}

module.exports = { ContractServiceImpl };

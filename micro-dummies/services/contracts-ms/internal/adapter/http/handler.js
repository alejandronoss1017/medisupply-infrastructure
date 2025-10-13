// HTTP handler - REST API endpoints
const express = require('express');

class ContractHandler {
    constructor(contractService) {
        this.contractService = contractService;
        this.router = express.Router();
        this.setupRoutes();
    }

    setupRoutes() {
        this.router.post('/contracts', this.createContract.bind(this));
        this.router.get('/contracts/:id', this.getContract.bind(this));
        this.router.put('/contracts/:id', this.updateContract.bind(this));
    }

    async createContract(req, res) {
        try {
            const result = await this.contractService.createContract(req.body);
            res.status(201).json({
                success: true,
                data: result,
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Error in createContract:', error);
            res.status(400).json({
                success: false,
                error: error.message,
                timestamp: new Date().toISOString()
            });
        }
    }

    async getContract(req, res) {
        try {
            const { id } = req.params;
            const result = await this.contractService.getContract(id);
            res.status(200).json({
                success: true,
                data: result,
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Error in getContract:', error);
            res.status(404).json({
                success: false,
                error: error.message,
                timestamp: new Date().toISOString()
            });
        }
    }

    async updateContract(req, res) {
        try {
            const { id } = req.params;
            const result = await this.contractService.updateContract(id, req.body);
            res.status(200).json({
                success: true,
                data: result,
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Error in updateContract:', error);
            res.status(400).json({
                success: false,
                error: error.message,
                timestamp: new Date().toISOString()
            });
        }
    }

    getRouter() {
        return this.router;
    }
}

module.exports = { ContractHandler };

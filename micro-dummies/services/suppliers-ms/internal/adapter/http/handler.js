// HTTP handler - REST API endpoints
const express = require('express');

class SupplierHandler {
    constructor(supplierService) {
        this.supplierService = supplierService;
        this.router = express.Router();
        this.setupRoutes();
    }

    setupRoutes() {
        this.router.post('/suppliers', this.createSupplier.bind(this));
        this.router.put('/suppliers/:id', this.updateSupplier.bind(this));
        this.router.post('/suppliers/process', this.processRequest.bind(this));
    }

    async createSupplier(req, res) {
        try {
            const result = await this.supplierService.createSupplier(req.body);
            res.status(201).json({
                success: true,
                data: result,
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Error in createSupplier:', error);
            res.status(400).json({
                success: false,
                error: error.message,
                timestamp: new Date().toISOString()
            });
        }
    }

    async updateSupplier(req, res) {
        try {
            const { id } = req.params;
            const result = await this.supplierService.updateSupplier(id, req.body);
            res.status(200).json({
                success: true,
                data: result,
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Error in updateSupplier:', error);
            res.status(400).json({
                success: false,
                error: error.message,
                timestamp: new Date().toISOString()
            });
        }
    }

    async processRequest(req, res) {
        try {
            const result = await this.supplierService.processSupplierRequest(req.body);
            res.status(200).json({
                success: true,
                data: result,
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Error in processRequest:', error);
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

module.exports = { SupplierHandler };

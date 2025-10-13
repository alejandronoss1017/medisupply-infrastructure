// HTTP handler - REST API endpoints
const express = require('express');

class PurchasePlanHandler {
    constructor(purchasePlanService) {
        this.purchasePlanService = purchasePlanService;
        this.router = express.Router();
        this.setupRoutes();
    }

    setupRoutes() {
        this.router.post('/purchase-plans', this.createPurchasePlan.bind(this));
        this.router.get('/purchase-plans/:id', this.getPurchasePlan.bind(this));
        this.router.put('/purchase-plans/:id', this.updatePurchasePlan.bind(this));
    }

    async createPurchasePlan(req, res) {
        try {
            const result = await this.purchasePlanService.createPurchasePlan(req.body);
            res.status(201).json({
                success: true,
                data: result,
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Error in createPurchasePlan:', error);
            res.status(400).json({
                success: false,
                error: error.message,
                timestamp: new Date().toISOString()
            });
        }
    }

    async getPurchasePlan(req, res) {
        try {
            const { id } = req.params;
            const result = await this.purchasePlanService.getPurchasePlan(id);
            res.status(200).json({
                success: true,
                data: result,
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Error in getPurchasePlan:', error);
            res.status(404).json({
                success: false,
                error: error.message,
                timestamp: new Date().toISOString()
            });
        }
    }

    async updatePurchasePlan(req, res) {
        try {
            const { id } = req.params;
            const result = await this.purchasePlanService.updatePurchasePlan(id, req.body);
            res.status(200).json({
                success: true,
                data: result,
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Error in updatePurchasePlan:', error);
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

module.exports = { PurchasePlanHandler };

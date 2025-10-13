// DynamoDB repository adapter
const { DynamoDBClient } = require('@aws-sdk/client-dynamodb');
const { DynamoDBDocumentClient, PutCommand, GetCommand, ScanCommand, UpdateCommand, DeleteCommand } = require('@aws-sdk/lib-dynamodb');
const { PurchasePlanRepository } = require('../../core/port/driven/repository');

class DynamoDBPurchasePlanRepository extends PurchasePlanRepository {
    constructor(config) {
        super();
        this.tableName = config.tableName || 'purchase-plans-db';

        const client = new DynamoDBClient({
            region: config.region || process.env.AWS_REGION || 'us-east-1',
            credentials: (config.accessKeyId && config.secretAccessKey) ? {
                accessKeyId: config.accessKeyId,
                secretAccessKey: config.secretAccessKey,
                sessionToken: config.sessionToken
            } : undefined,
        });

        this.ddb = DynamoDBDocumentClient.from(client);
        console.log(`DynamoDB repository initialized with table: ${this.tableName}`);
    }

    async save(purchasePlan) {
        try {
            const item = {
                ...purchasePlan,
                ts: new Date().toISOString()
            };

            await this.ddb.send(new PutCommand({
                TableName: this.tableName,
                Item: item
            }));

            console.log(`Purchase plan saved to DynamoDB: ${purchasePlan.id}`);
            return item;
        } catch (error) {
            console.error('Error saving to DynamoDB:', error);
            throw error;
        }
    }

    async findById(id) {
        try {
            const result = await this.ddb.send(new GetCommand({
                TableName: this.tableName,
                Key: { id }
            }));

            return result.Item || null;
        } catch (error) {
            console.error('Error finding by ID in DynamoDB:', error);
            throw error;
        }
    }

    async findAll() {
        try {
            const result = await this.ddb.send(new ScanCommand({
                TableName: this.tableName
            }));

            return result.Items || [];
        } catch (error) {
            console.error('Error scanning DynamoDB:', error);
            throw error;
        }
    }

    async update(id, purchasePlan) {
        try {
            const item = {
                ...purchasePlan,
                id,
                updatedAt: new Date().toISOString()
            };

            await this.ddb.send(new PutCommand({
                TableName: this.tableName,
                Item: item
            }));

            return item;
        } catch (error) {
            console.error('Error updating in DynamoDB:', error);
            throw error;
        }
    }

    async delete(id) {
        try {
            await this.ddb.send(new DeleteCommand({
                TableName: this.tableName,
                Key: { id }
            }));

            console.log(`Purchase plan deleted from DynamoDB: ${id}`);
            return true;
        } catch (error) {
            console.error('Error deleting from DynamoDB:', error);
            throw error;
        }
    }
}

module.exports = { DynamoDBPurchasePlanRepository };

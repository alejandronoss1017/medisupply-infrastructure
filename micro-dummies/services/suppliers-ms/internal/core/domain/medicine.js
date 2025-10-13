// Domain model for medicine
class Medicine {
    constructor({ id, name, description, price, category, supplierId = {} }) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.price = price;
        this.category = category;
        this.supplierId = supplierId;
        this.createdAt = new Date();
        this.updatedAt = new Date();
    }

    update(data) {
        Object.keys(data).forEach(key => {
            if (this.hasOwnProperty(key) && key !== 'id' && key !== 'createdAt') {
                    this[key] = data[key];
            }
        });
        this.updatedAt = new Date();
    }

    validate() {
        if (!this.name || this.name.trim() === '') {
            throw new Error('Medicine name is required');
        }
        return true;
    }

    toJSON() {
        return {
            id: this.id,
            name: this.name,
            description: this.description,
            price: this.price,
            category: this.category,
            supplierId: this.supplierId,
            createdAt: this.createdAt,
            updatedAt: this.updatedAt
        };
    }
}

module.exports = { Medicine };

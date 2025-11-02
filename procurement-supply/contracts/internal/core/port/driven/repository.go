package driven

// Repository defines a generic interface for data persistence
// ID is the type of the entity's identifier (must be comparable for map keys)
// T is the entity type
type Repository[ID comparable, T any] interface {
	// Create adds a new entity to the repository
	Create(entity T) error

	// FindByID retrieves an entity by its ID
	FindByID(id ID) (*T, error)

	// FindAll retrieves all entities
	FindAll() ([]T, error)

	// Update modifies an existing entity
	Update(entity T) error

	// Delete removes an entity by its ID
	Delete(id ID) error

	// Exists checks if an entity with the given ID exists
	Exists(id ID) bool
}

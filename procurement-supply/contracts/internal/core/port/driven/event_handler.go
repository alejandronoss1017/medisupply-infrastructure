package driven

import "contracts/internal/core/domain"

// MedicineEventHandler represents the driven port for handling medicine events
// This interface defines what external services our application needs
type MedicineEventHandler interface {
	HandleMedicineCreated(event *domain.Event[domain.Medicine]) error
	HandleMedicineUpdated(event *domain.Event[domain.Medicine]) error
	HandleMedicineDeleted(event *domain.Event[domain.Medicine]) error
}

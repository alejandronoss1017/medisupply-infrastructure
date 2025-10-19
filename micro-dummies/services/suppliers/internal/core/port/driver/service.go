package driver

import "suppliers/internal/core/domain"

type MedicineService interface {
	RetrieveMedicines() ([]domain.Medicine, error)
	RetrieveMedicine(id string) (*domain.Medicine, error)
	CreateMedicine(medicine *domain.Medicine) (*domain.Medicine, error)
	UpdateMedicine(id string, medicine *domain.Medicine) (*domain.Medicine, error)
	DeleteMedicine(id string) error
}

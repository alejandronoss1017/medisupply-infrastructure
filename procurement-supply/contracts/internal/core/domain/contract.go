package domain

type Contract struct {
	ID         string `json:"id"`
	Path       string `json:"path"`
	CustomerID string `json:"customerId"`
	SLAs       []SLA  `json:"slas"`
}

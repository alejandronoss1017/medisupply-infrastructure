package domain

type Contract struct {
	ID         string `json:"id" abi:"id"`
	Path       string `json:"path" abi:"path"`
	CustomerID string `json:"customerId" abi:"clientId"`
	SLAs       []*SLA `json:"slas" abi:"slas"`
}

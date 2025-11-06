package domain

import "math/big"

type Novelty struct {
	ID         string   `json:"id,omitempty"`
	ContractID string   `json:"contractId"`
	CustomerID string   `json:"customerId"`
	SLAID      string   `json:"slaId"`
	Value      *big.Int `json:"value"`
}

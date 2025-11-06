package domain

import "math/big"

type Comparator uint8
type SLAStatus uint8

const (
	GreaterThan Comparator = iota
	LessThan
	Equal
	GreaterOrEqual
	LessOrEqual
)

const (
	Active SLAStatus = iota
	Violated
	Compliant
	Inactive
)

type SLA struct {
	ID          string     `json:"id" abi:"id"`
	Name        string     `json:"name" abi:"name"`
	Description string     `json:"description" abi:"description"`
	Target      *big.Int   `json:"target" abi:"target"`
	Comparator  Comparator `json:"comparator" abi:"comparator"`
	Status      SLAStatus  `json:"status" abi:"status"`
}

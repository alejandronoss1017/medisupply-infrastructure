package domain

import (
	"math/big"
	"time"
)

// BlockchainEvent represents a generic blockchain event
type BlockchainEvent struct {
	EventType   EventType
	BlockNumber uint64
	TxHash      string
	Timestamp   time.Time
}

// EventType represents the type blockchain event
type EventType string

const (
	EventTypeContractAdded    EventType = "ContractAdded"
	EventTypeSLAAdded         EventType = "SLAAdded"
	EventTypeSLAStatusUpdated EventType = "SLAStatusUpdated"
	SLAViolated               EventType = "SLAViolated"
)

// ContractAddedEvent represents the ContractAdded event from the blockchain
type ContractAddedEvent struct {
	BlockchainEvent
	ContractID string
	CustomerID string
}

// SLAAddedEvent represents the SLAAdded event from the blockchain
type SLAAddedEvent struct {
	BlockchainEvent
	ContractID string
	SLAID      string
}

// SLAStatusUpdatedEvent represents the SLAStatusUpdated event from the blockchain
type SLAStatusUpdatedEvent struct {
	BlockchainEvent
	ContractID string
	SLAID      string
	NewStatus  uint8
}

// SLAStatus represents the status of an SLA
type SLAStatus uint8

const (
	ACTIVE    SLAStatus = 0
	VIOLATED  SLAStatus = 1
	COMPLIANT SLAStatus = 2
	INACTIVE  SLAStatus = 3
)

// String returns the string representation of SLAStatus
func (s SLAStatus) String() string {
	switch s {
	case ACTIVE:
		return "Active"
	case VIOLATED:
		return "Violated"
	case COMPLIANT:
		return "Compliant"
	case INACTIVE:
		return "Inactive"
	default:
		return "Unknown"
	}
}

// Comparator represents the comparison operator for SLA
type Comparator uint8

const (
	GREATER_THAN     Comparator = 0
	LESS_THAN        Comparator = 1
	EQUAL            Comparator = 2
	GREATER_OR_EQUAL Comparator = 3
	LESS_OR_EQUAL    Comparator = 4
)

// String returns the string representation of Comparator
func (c Comparator) String() string {
	switch c {
	case GREATER_THAN:
		return "GreaterThan"
	case LESS_THAN:
		return "LessThan"
	case EQUAL:
		return "Equal"
	case GREATER_OR_EQUAL:
		return "GreaterOrEqual"
	case LESS_OR_EQUAL:
		return "LessOrEqualThan"
	}
	return "Unknown"
}

// SLA represents a Service Level Agreement
type SLA struct {
	ID          string
	Name        string
	Description string
	Target      *big.Int
	Comparator  Comparator
	Status      SLAStatus
}

// Contract represents a smart contract with its associated SLAs
type Contract struct {
	ID         string
	Path       string
	CustomerID string
	SLAs       []*SLA
}

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
	SLAStatusPending  SLAStatus = 0
	SLAStatusMet      SLAStatus = 1
	SLAStatusViolated SLAStatus = 2
)

// String returns the string representation of SLAStatus
func (s SLAStatus) String() string {
	switch s {
	case SLAStatusPending:
		return "Pending"
	case SLAStatusMet:
		return "Met"
	case SLAStatusViolated:
		return "Violated"
	default:
		return "Unknown"
	}
}

// Comparator represents the comparison operator for SLA
type Comparator uint8

const (
	ComparatorLessThan    Comparator = 0
	ComparatorGreaterThan Comparator = 1
	ComparatorEqualTo     Comparator = 2
)

// String returns the string representation of Comparator
func (c Comparator) String() string {
	switch c {
	case ComparatorLessThan:
		return "LessThan"
	case ComparatorGreaterThan:
		return "GreaterThan"
	case ComparatorEqualTo:
		return "EqualTo"
	default:
		return "Unknown"
	}
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

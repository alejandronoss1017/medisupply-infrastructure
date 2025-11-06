package logger

import (
	"time"

	"go.uber.org/zap"
)

// Common structured field helpers for consistent logging across the application

// String creates a string field
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int creates an integer field
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Int64 creates an int64 field
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Uint64 creates a uint64 field
func Uint64(key string, val uint64) zap.Field {
	return zap.Uint64(key, val)
}

// Bool creates a boolean field
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Duration creates a duration field
func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

// Error creates an error field
func Error(err error) zap.Field {
	return zap.Error(err)
}

// Any creates a field with any type (uses reflection)
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

// ContractID is a convenience helper for contract IDs
func ContractID(id string) zap.Field {
	return zap.String("contract_id", id)
}

// CustomerID is a convenience helper for customer IDs
func CustomerID(id string) zap.Field {
	return zap.String("customer_id", id)
}

// SLAID is a convenience helper for SLA IDs
func SLAID(id string) zap.Field {
	return zap.String("sla_id", id)
}

// TxHash is a convenience helper for transaction hashes
func TxHash(hash string) zap.Field {
	return zap.String("tx_hash", hash)
}

// BlockNumber is a convenience helper for block numbers
func BlockNumber(num uint64) zap.Field {
	return zap.Uint64("block_number", num)
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package dbgen

import (
	"database/sql/driver"
	"fmt"
)

type PetsStatus string

const (
	PetsStatusAvailable PetsStatus = "available"
	PetsStatusPending   PetsStatus = "pending"
	PetsStatusSold      PetsStatus = "sold"
)

func (e *PetsStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PetsStatus(s)
	case string:
		*e = PetsStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PetsStatus: %T", src)
	}
	return nil
}

type NullPetsStatus struct {
	PetsStatus PetsStatus
	Valid      bool // Valid is true if PetsStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPetsStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PetsStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PetsStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPetsStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PetsStatus), nil
}

type Pet struct {
	ID     int64
	Name   string
	Status PetsStatus
}

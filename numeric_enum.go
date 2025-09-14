package enums

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// Int is an interned int.
type Int int

// ClosedNumericEnum is numeric enum that is saved in db on in json as a readable string representation.
type ClosedNumericEnum interface {
	driver.Valuer
	sql.Scanner
	json.Marshaler
	json.Unmarshaler
	EnumValueIsValid() bool
	DefaultValue() int
}

// NumericEnum generic type for numeric closed enum.
type NumericEnum[T ~int] struct {
	values        map[T]string
	reveredValues map[string]T
}

// NewNumericClosedEnum creates a new closed numeric enum.
func NewNumericClosedEnum[T ~int](values map[T]string) NumericEnum[T] {
	reveredValues := make(map[string]T, len(values))
	for v, s := range values {
		reveredValues[s] = v
	}

	if _, found := values[0]; !found {
		panic("numeric enum must have default value 0")
	}

	return NumericEnum[T]{
		values:        values,
		reveredValues: reveredValues,
	}
}

// DefaultValue returns default value for numeric enum.
func (e *NumericEnum[T]) DefaultValue() T {
	return 0
}

// GetStringRepresentation returns string representation of enum value.
func (e *NumericEnum[T]) GetStringRepresentation(v T) (string, bool) {
	s, f := e.values[v]
	return s, f
}

// Get returns an enum value for provided int representation.
func (e *NumericEnum[T]) Get(v int) (T, bool) {
	if _, found := e.values[T(v)]; found {
		return T(v), true
	}
	return 0, false
}

// GetFromString returns an enum value for provided string representation.
func (e *NumericEnum[T]) GetFromString(s string) (T, bool) {
	v, found := e.reveredValues[s]
	return v, found
}

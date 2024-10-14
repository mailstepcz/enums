package enums

import (
	"database/sql/driver"
	"errors"
	"strings"
	"unsafe"

	"github.com/fealsamh/datastructures/unionfind"
)

// ErrTransitionNotAllowed signifies a forbidden transition.
var ErrTransitionNotAllowed = errors.New("transition not allowed")

// String is an interned string.
type String string

// ClosedEnum is a closed enum.
type ClosedEnum interface {
	driver.Valuer
	EnumValueIsValid() bool
	DefaultValue() string
}

// Transitioner is a transition verifier.
type Transitioner[T ClosedEnum] interface {
	CanTransition(T) bool
	Eq(T) bool
}

// Compare compares two interned strings.
func (s1 String) Compare(s2 String) int {
	return strings.Compare(string(s1), string(s2))
}

// Eq checks two interned strings for equality.
func (s1 String) Eq(s2 String) bool {
	return unsafe.StringData(string(s1)) == unsafe.StringData(string(s2))
}

// Enum is a universal enum type.
type Enum struct {
	values       *unionfind.Structure[String]
	defaultValue String
}

// DefaultValue returns the default value of the enum.
func (e *Enum) DefaultValue() String {
	return e.defaultValue
}

// EnumGet returns an interned enum value for the provided string.
func EnumGet[T interface {
	~string
	ClosedEnum
}](e *Enum, s string) (T, bool) {
	r, ok := e.Get(s)
	return T(r), ok
}

// NewClosedEnum creates a new closed enum.
func NewClosedEnum[T ~string](values ...T) Enum {
	uf := unionfind.New[String]()
	var defaultValue String
	for i, s := range values {
		s, _ := uf.Add(String(s))
		if i == 0 {
			defaultValue = s.Value
		}
	}
	return Enum{values: uf, defaultValue: defaultValue}
}

// Get returns an interned string for the provided enum value.
func (e *Enum) Get(s string) (String, bool) {
	r, ok := e.values.Get(String(s))
	if ok {
		return r.Value, true
	}
	return "", false
}

// Transition is a state transition.
type Transition[T ClosedEnum] struct {
	From T
	To   T
}

// AllowedTransitions returns a slice of allowed transitions.
func AllowedTransitions[E interface {
	ClosedEnum
	Transitioner[E]
}](states ...E) ([]Transition[E], error) {
	var ts []Transition[E]
	for _, from := range states {
		for _, to := range states {
			if from.CanTransition(to) {
				if from.Eq(to) {
					return nil, ErrTransitionNotAllowed
				}
				ts = append(ts, Transition[E]{
					From: from,
					To:   to,
				})
			}
		}
	}
	return ts, nil
}

package enums

import (
	"database/sql/driver"
	"strings"
	"unsafe"

	"github.com/fealsamh/datastructures/unionfind"
)

// String is an interned string.
type String string

// ClosedEnum is a closed enum.
type ClosedEnum interface {
	driver.Valuer
	EnumValueIsValid() bool
	DefaultValue() string
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

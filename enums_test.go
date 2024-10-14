package enums

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

type AbcEnum String

const (
	A1 AbcEnum = "a1"
	B2 AbcEnum = "b2"
	C3 AbcEnum = "c3"
)

var abcEnum = NewClosedEnum(A1, B2, C3)

var (
	_ ClosedEnum            = AbcEnum("")
	_ Transitioner[AbcEnum] = AbcEnum("")
)

func (v AbcEnum) EnumValueIsValid() bool {
	_, ok := EnumGet[AbcEnum](&abcEnum, string(v))
	return ok
}

func (v AbcEnum) DefaultValue() string {
	return string(abcEnum.DefaultValue())
}

func (v AbcEnum) Value() (driver.Value, error) {
	return string(v), nil
}

func (v AbcEnum) CanTransition(newState AbcEnum) bool {
	switch v {
	case A1:
		return newState == B2
	case B2:
		return newState == C3
	}
	return false
}

func TestClosedEnum(t *testing.T) {
	r := require.New(t)

	s, ok := EnumGet[AbcEnum](&abcEnum, "a1")
	r.True(ok)
	r.Equal("a1", string(s))

	s, ok = EnumGet[AbcEnum](&abcEnum, "d4")
	r.False(ok)
	r.Equal("", string(s))
}

func TestEnumInterned(t *testing.T) {
	r := require.New(t)

	s1 := fmt.Sprintf("a%d", 1)
	s2 := fmt.Sprintf("a%d", 1)

	r.True(unsafe.StringData(s1) != unsafe.StringData(s2))

	e1, ok := EnumGet[AbcEnum](&abcEnum, s1)
	r.True(ok)
	e2, ok := EnumGet[AbcEnum](&abcEnum, s2)
	r.True(ok)

	r.True(unsafe.StringData(string(e1)) == unsafe.StringData(string(e2)))
}

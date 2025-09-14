package enums

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/mailstepcz/serr"
	"github.com/stretchr/testify/require"
)

type SomeEnum Int

var (
	_ ClosedNumericEnum = (*SomeEnum)(nil)
)

const (
	SomeEnum0 SomeEnum = iota
	SomeEnum1
)

var (
	enumStringMap = map[SomeEnum]string{
		SomeEnum0: "Value0",
		SomeEnum1: "Value1",
	}
)

type Person struct {
	Age SomeEnum `json:"age"`
}

var someEnum = NewNumericClosedEnum(enumStringMap)

func (se SomeEnum) DefaultValue() int {
	return int(someEnum.DefaultValue())
}

func (se SomeEnum) EnumValueIsValid() bool {
	_, found := someEnum.Get(int(se))
	return found
}

func (se SomeEnum) MarshalJSON() ([]byte, error) {
	s, f := someEnum.GetStringRepresentation(se)
	if !f {
		return nil, serr.New("invalid enum value", serr.Int("enumValue", int(se)))
	}
	return json.Marshal(s)
}

func (se *SomeEnum) UnmarshalJSON(v []byte) error {
	var raw string
	if err := json.Unmarshal(v, &raw); err != nil {
		return err
	}
	e, f := someEnum.GetFromString(raw)
	if !f {
		return serr.New("invalid enum value", serr.String("enumValue", string(v)))
	}
	*se = e
	return nil
}

func (se *SomeEnum) Scan(src any) error {
	return nil
}

func (se SomeEnum) Value() (driver.Value, error) {
	s, f := someEnum.GetStringRepresentation(se)
	if !f {
		return nil, serr.New("invalid enum value", serr.Int("enumValue", int(se)))
	}
	return s, nil
}

func TestNumericEnumGetDefaultValue(t *testing.T) {
	req := require.New(t)

	v := someEnum.DefaultValue()
	req.Equal(SomeEnum0, v)
}

func TestMarshalNumericEnum(t *testing.T) {
	req := require.New(t)
	p := Person{
		Age: SomeEnum1,
	}

	b, err := json.Marshal(&p)
	req.NoError(err)
	req.Equal("{\"age\":\"Value1\"}", string(b))
}

func TestUnmarshalNumericEnum(t *testing.T) {
	req := require.New(t)
	var p Person
	err := json.Unmarshal([]byte("{\"age\":\"Value1\"}"), &p)
	req.NoError(err)
	req.Equal(p.Age, SomeEnum1)
}

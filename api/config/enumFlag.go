package config

import (
	"fmt"
	"strings"
)

type EnumValue struct {
	allowedValues []string
	value         string
}

func (e *EnumValue) String() string {
	return e.value
}

func (e *EnumValue) Set(v string) error {
	for _, allowed := range e.allowedValues {
		if v == allowed {
			e.value = v
			return nil
		}
	}
	return fmt.Errorf("invalid value %q, must be one of: %s", v, strings.Join(e.allowedValues, ", "))
}

func NewEnumValue(defaultValue string, allowedValues []string) *EnumValue {
	return &EnumValue{
		allowedValues: allowedValues,
		value:         defaultValue,
	}
}

func (e *EnumValue) AllowedValuesList() string {
	return fmt.Sprintf("(one of: %s)", strings.Join(e.allowedValues, ", "))
}

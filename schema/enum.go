package schema

import (
	"errors"
	"fmt"
	"strings"
)

// EnumDefinition represents the definition of an enum type.
type EnumDefinition struct {
	// Values is a list of distinct, non-empty values that are part of the enum type.
	// Each value must conform to the NameFormat regular expression.
	Values []string
}

func (EnumDefinition) isType() {}

// Validate validates the enum definition.
func (e EnumDefinition) Validate() error {
	if len(e.Values) == 0 {
		return errors.New("enum definition values cannot be empty")
	}
	seen := make(map[string]bool, len(e.Values))
	for i, v := range e.Values {
		if !ValidateName(v) {
			return fmt.Errorf("invalid enum definition value %q at index %d", v, i)
		}

		if seen[v] {
			return fmt.Errorf("duplicate enum definition value %q", v)
		}
		seen[v] = true
	}
	return nil
}

func (e EnumDefinition) validateWithSchema(map[string]Type) error {
	return e.Validate()
}

// ValidateValue validates that the value is a valid enum value.
func (e EnumDefinition) ValidateValue(value string) error {
	for _, v := range e.Values {
		if v == value {
			return nil
		}
	}
	return fmt.Errorf("value %q is not a valid enum value, must be one of: %s", value, strings.Join(e.Values, ", "))
}

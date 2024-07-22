package schema

import "fmt"

// ModuleSchema represents the logical schema of a module for purposes of indexing and querying.
type ModuleSchema struct {
	Types map[string]Type
}

// Validate validates the module schema.
func (s ModuleSchema) Validate() error {
	for _, typ := range s.Types {
		if err := typ.validateWithSchema(s.Types); err != nil {
			return err
		}
	}

	return nil
}

// ValidateObjectUpdate validates that the update conforms to the module schema.
func (s ModuleSchema) ValidateObjectUpdate(update ObjectUpdate) error {
	typ, ok := s.Types[update.TypeName]
	if !ok {
		return fmt.Errorf("object type %q not found in module schema", update.TypeName)
	}

	objType, ok := typ.(ObjectType)
	if !ok {
		return fmt.Errorf("object type %q is not an object type", update.TypeName)
	}

	return objType.ValidateObjectUpdate(update)
}

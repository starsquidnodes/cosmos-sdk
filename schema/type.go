package schema

type Type interface {
	isType()
	Validate() error
	validateWithSchema(types map[string]Type) error
}

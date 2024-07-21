package schema

type Type interface {
	Name() string
	Validate() error

	isType()
	validateWithSchema(types map[string]Type) error
}

package schema

type Type interface {
	TypeName() string
	Validate() error

	isType()
	validateWithSchema(types map[string]Type) error
}

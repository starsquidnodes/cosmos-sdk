package indexertesting

import (
	"fmt"

	"cosmossdk.io/schema"
)

var ExampleAppSchema = map[string]schema.ModuleSchema{
	"all_kinds": mkAllKindsModule(),
	"test_cases": {
		ObjectTypes: []schema.ObjectType{
			{
				"Singleton",
				[]schema.Field{},
				[]schema.Field{
					{
						Name: "Value",
						Kind: schema.StringKind,
					},
					{
						Name: "Value2",
						Kind: schema.BytesKind,
					},
				},
				false,
			},
			{
				Name: "Simple",
				KeyFields: []schema.Field{
					{
						Name: "Key",
						Kind: schema.StringKind,
					},
				},
				ValueFields: []schema.Field{
					{
						Name: "Value1",
						Kind: schema.Int32Kind,
					},
					{
						Name: "Value2",
						Kind: schema.BytesKind,
					},
				},
			},
			{
				Name: "TwoKeys",
				KeyFields: []schema.Field{
					{
						Name: "Key1",
						Kind: schema.StringKind,
					},
					{
						Name: "Key2",
						Kind: schema.Int32Kind,
					},
				},
			},
			{
				Name: "ThreeKeys",
				KeyFields: []schema.Field{
					{
						Name: "Key1",
						Kind: schema.StringKind,
					},
					{
						Name: "Key2",
						Kind: schema.Int32Kind,
					},
					{
						Name: "Key3",
						Kind: schema.Uint64Kind,
					},
				},
				ValueFields: []schema.Field{
					{
						Name: "Value1",
						Kind: schema.Int32Kind,
					},
				},
			},
			{
				Name: "ManyValues",
				KeyFields: []schema.Field{
					{
						Name: "Key",
						Kind: schema.StringKind,
					},
				},
				ValueFields: []schema.Field{
					{
						Name: "Value1",
						Kind: schema.Int32Kind,
					},
					{
						Name: "Value2",
						Kind: schema.BytesKind,
					},
					{
						Name: "Value3",
						Kind: schema.Float64Kind,
					},
					{
						Name: "Value4",
						Kind: schema.Uint64Kind,
					},
				},
			},
			{
				Name: "RetainDeletions",
				KeyFields: []schema.Field{
					{
						Name: "Key",
						Kind: schema.StringKind,
					},
				},
				ValueFields: []schema.Field{
					{
						Name: "Value1",
						Kind: schema.Int32Kind,
					},
					{
						Name: "Value2",
						Kind: schema.BytesKind,
					},
				},
				RetainDeletions: true,
			},
		},
	},
}

func mkAllKindsModule() schema.ModuleSchema {
	mod := schema.ModuleSchema{}

	for i := 1; i < int(schema.MAX_VALID_KIND); i++ {
		kind := schema.Kind(i)
		typ := mkTestObjectType(kind)
		mod.ObjectTypes = append(mod.ObjectTypes, typ)
	}

	return mod
}

func mkTestObjectType(kind schema.Kind) schema.ObjectType {
	field := schema.Field{
		Kind: kind,
	}

	if kind == schema.EnumKind {
		field.EnumDefinition = testEnum
	}

	if kind == schema.Bech32AddressKind {
		field.AddressPrefix = "cosmos"
	}

	key1Field := field
	key1Field.Name = "keyNotNull"
	key2Field := field
	key2Field.Name = "keyNullable"
	key2Field.Nullable = true
	val1Field := field
	val1Field.Name = "valNotNull"
	val2Field := field
	val2Field.Name = "valNullable"
	val2Field.Nullable = true

	return schema.ObjectType{
		Name:        fmt.Sprintf("test_%v", kind),
		KeyFields:   []schema.Field{key1Field, key2Field},
		ValueFields: []schema.Field{val1Field, val2Field},
	}
}

var testEnum = schema.EnumDefinition{
	Name:   "test_enum",
	Values: []string{"foo", "bar", "baz"},
}
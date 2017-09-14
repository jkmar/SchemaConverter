package crud

import "fmt"

// GenerateList generates a list function
func GenerateList(
	goextPackage,
	name,
	typeName string,
	lock,
	raw bool,
) string {
	var (
		prefix,
		suffix,
		arg,
		argType string
	)
	if raw {
		suffix = "Raw"
	}
	if lock {
		prefix = "Lock"
		arg = ", policy"
		argType = " " + goextPackage + ".LockPolicy"
	}

	return fmt.Sprintf(
		`func %sList%s%s(`+
			`schema %s.ISchema, `+
			`filter %s.Filter, `+
			`paginator *%s.Paginator, `+
			`context %s.Context%s%s) ([]%s, error) {
	list, err := schema.%sList%s(filter, paginator, context%s)
	if err != nil {
		return nil, err
	}
	result := make([]%s, len(list))
	for i, object := range list {
		result[i] = object.(%s)
	}
	return result, nil
}
`,
		prefix,
		suffix,
		name,
		goextPackage,
		goextPackage,
		goextPackage,
		goextPackage,
		arg,
		argType,
		typeName,
		prefix,
		suffix,
		arg,
		typeName,
		typeName,
	)
}

// GenerateFetch generates a fetch function
func GenerateFetch(
	goextPackage,
	name,
	typeName string,
	lock,
	raw bool,
) string {
	var (
		prefix,
		suffix,
		arg,
		argType string
	)
	if raw {
		suffix = "Raw"
	}
	if lock {
		prefix = "Lock"
		arg = ", policy"
		argType = " " + goextPackage + ".LockPolicy"
	}

	return fmt.Sprintf(
		`func %sFetch%s%s(`+
			`schema %s.ISchema, `+
			`id string, `+
			`context %s.Context%s%s) (%s, error) {
	result, err := schema.%sFetch%s(id, context%s)
	if err != nil {
		return nil, err
	}
	return result.(%s), nil
}
`,
		prefix,
		suffix,
		name,
		goextPackage,
		goextPackage,
		arg,
		argType,
		typeName,
		prefix,
		suffix,
		arg,
		typeName,
	)
}

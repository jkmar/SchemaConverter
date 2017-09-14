package crud

import "fmt"

// GenerateList generates a list function
func GenerateList(
	goextPackage,
	name,
	typeName,
	prefix,
	suffix string,
) string {
	return fmt.Sprintf(
		`func %sList%s%s(schema %s.ISchema, filter %s.Filter, paginator *%s.Paginator, context %s.Context) ([]%s, error) {
	list, err := schema.%sList%s(filter, paginator, context)
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
		typeName,
		prefix,
		suffix,
		typeName,
		typeName,
	)
}

// GenerateFetch generates a fetch function
func GenerateFetch(
	goextPackage,
	name,
	typeName,
	prefix,
	suffix string,
) string {
	return fmt.Sprintf(
		`func %sFetch%s%s(schema %s.ISchema, id string, context %s.Context) (%s, error) {
	result, err := schema.%sFetch%s(id, context)
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
		typeName,
		prefix,
		suffix,
		typeName,
	)
}

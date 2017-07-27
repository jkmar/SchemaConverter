package main

type Item interface {
	Type() string
	Parse(string, map[interface{}]interface{}) error
}

func CreateItem(itemType interface{}) (Item, error) {
	strType, err := parseType(itemType)
	if err != nil {
		return nil, err
	}
	return CreateItemFromString(strType), nil
}

func CreateItemFromString(itemType string) Item {
	switch itemType {
	case "array":
		return &Array{}
	case "object":
		return &Object{}
	default:
		return &PlainItem{}
	}
}

package main

type PlainItem struct {
	ItemType string
}

func (item *PlainItem) Type() string {
	return item.ItemType
}

func (item *PlainItem) Parse(prefix string, object map[interface{}]interface{}) {
	item.ItemType = parseType(object["type"])
}

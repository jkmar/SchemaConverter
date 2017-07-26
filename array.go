package main

type Array struct {
	Item Item
}

func (item *Array) Type() string {
	return "[]" + item.Item.Type()
}

func (item *Array) Parse(prefix string, object map[interface{}]interface{}) {
	next := object["items"].(map[interface{}]interface{})
	item.Item = CreateItem(next["type"])
	item.Item.Parse(prefix, next)
}

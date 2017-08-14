package app

import (
	"github.com/zimnx/YamlSchemaToGoStruct/reader"
	"github.com/zimnx/YamlSchemaToGoStruct/schema"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
	"github.com/zimnx/YamlSchemaToGoStruct/writer"
)

func readConfig(config, input string) ([]map[interface{}]interface{}, error) {
	if len(config) == 0 {
		return nil, nil
	}
	return reader.ReadAll(config, input)
}

func writeResult(data []string, input, output, suffix string) error {
	rawData := util.CollectData(input, data)
	if rawData == "" {
		return nil
	}
	file := writer.CreateWriter(util.TryToAddName(output, suffix))
	return file.Write(rawData)
}

// Run application
func Run(input, output, config, suffix string) error {
	other, err := readConfig(config, input)
	if err != nil {
		return err
	}
	objects, err := reader.ReadSingle(input)
	if err != nil {
		return err
	}
	interfaces, structs, implementations, err := schema.Convert(other, objects, suffix)
	if err != nil {
		return err
	}
	if err = writeResult(interfaces, "interface", output, "interface.go"); err != nil {
		return err
	}
	if err = writeResult(structs, input, output, "raw.go"); err != nil {
		return err
	}
	return writeResult(implementations, input, output, "implementation.go")
}

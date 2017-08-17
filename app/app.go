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

func writeResult(data []string, packageName, outputPrefix, outputSuffix string) error {
	rawData := util.CollectData(packageName, data)
	if rawData == "" {
		return nil
	}
	file := writer.CreateWriter(util.TryToAddName(outputPrefix, outputSuffix))
	return file.Write(rawData)
}

// Run application
func Run(config, output, packageName, suffix string) error {
	all, err := readConfig(config, "")
	if err != nil {
		return err
	}

	interfaces, structs, implementations, err := schema.Convert(nil, all, suffix)
	if err != nil {
		return err
	}
	if err = writeResult(interfaces, "esi", output, "interface.go"); err != nil {
		return err
	}
	if err = writeResult(structs, packageName, output, "raw.go"); err != nil {
		return err
	}
	return writeResult(implementations, packageName, output, "implementation.go")
}

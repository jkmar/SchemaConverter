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
func Run(
	config,
	output,
	goextPackage,
	goodiesPackage,
	resourcePackage,
	interfacePackage,
	rawSuffix,
	interfaceSuffix string,
) error {
	interfaceSuffix = util.AddName(rawSuffix, interfaceSuffix)
	all, err := readConfig(config, "")
	if err != nil {
		return err
	}

	generated, err := schema.Convert(
		nil,
		all,
		rawSuffix,
		interfaceSuffix,
		goextPackage,
	)
	if err != nil {
		return err
	}

	if err = writeResult(
		generated.RawInterfaces,
		interfacePackage,
		output,
		"generated_interface.go",
	); err != nil {
		return err
	}

	if err = writeResult(
		generated.Interfaces,
		interfacePackage,
		output,
		"interface.go",
	); err != nil {
		return err
	}

	if err = writeResult(
		generated.Structs,
		resourcePackage,
		output,
		"raw.go",
	); err != nil {
		return err
	}

	if err = writeResult(
		generated.Implementations,
		resourcePackage,
		output,
		"implementation.go",
	); err != nil {
		return err
	}

	if err = writeResult(
		generated.Constructors,
		resourcePackage,
		output,
		"constructors.go",
	); err != nil {
		return err
	}

	if err = writeResult(
		generated.Crud,
		goodiesPackage,
		output,
		"crud.go",
	); err != nil {
		return err
	}

	return writeResult(
		generated.RawCrud,
		goodiesPackage,
		output,
		"raw_crud.go",
	)
}

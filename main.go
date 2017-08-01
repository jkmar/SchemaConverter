package main

import (
	"flag"
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/reader"
	"github.com/zimnx/YamlSchemaToGoStruct/schema"
)

func main() {
	config := flag.String("config", "", "config file")
	inputSchema := flag.String("schema", "", "path to yaml with schema")
	suffix := flag.String("suffix", "", "suffix appended to struct names")
	flag.Parse()

	other, err := reader.ReadAll(*config, *inputSchema)
	if err != nil {
		panic(err)
	}

	objects, err := reader.ReadSingle(*inputSchema)
	if err != nil {
		panic(err)
	}

	tmp, err := schema.Convert(other, objects, "db", "json", *suffix)
	if err != nil {
		panic(err)
	}
	for _, x := range tmp {
		fmt.Println(x)
	}
}

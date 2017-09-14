package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zimnx/YamlSchemaToGoStruct/app"
	"os"
)

var (
	goextPackage,
	goodiesPackage,
	resourcePackage,
	interfacePackage,
	output,
	rawSuffix,
	interfaceSuffix string
)

// RootCmd of application
var RootCmd = &cobra.Command{
	Use:   "YamlSchemaToGoStruct [path to config file with schemas]",
	Short: "YamlSchemaToGoStruct generates go structs from yaml schemas",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := app.Run(
			args[0],
			output,
			goextPackage,
			goodiesPackage,
			resourcePackage,
			interfacePackage,
			rawSuffix,
			interfaceSuffix,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.Flags().StringVar(
		&goextPackage,
		"goext",
		"goext",
		"package name for golang extension interfaces",
	)
	RootCmd.Flags().StringVar(
		&goodiesPackage,
		"goodies",
		"goodies",
		"package name for crud packages",
	)
	RootCmd.Flags().StringVarP(
		&resourcePackage,
		"package-name",
		"p",
		"resources",
		"package name for raw structs",
	)
	RootCmd.Flags().StringVarP(
		&interfacePackage,
		"interface-name",
		"i",
		"esi",
		"package name for interfaces",
	)
	RootCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"",
		"prefix added to output file",
	)
	RootCmd.Flags().StringVar(
		&rawSuffix,
		"raw-suffix",
		"",
		"suffix added to raw struct names",
	)
	RootCmd.Flags().StringVar(
		&interfaceSuffix,
		"interface-suffix",
		"gen",
		"suffix added to generated interface names",
	)
}

// Execute RootCmd
func Execute() {
	RootCmd.Execute()
}

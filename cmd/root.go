package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zimnx/YamlSchemaToGoStruct/app"
	"fmt"
	"os"
)

var (
	config         string
	annotationDB   string
	annotationJSON string
	suffix         string
	output         string
)

var RootCmd = &cobra.Command{
	Use:   "YamlSchemaToGoStruct [path to schema]",
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
			config,
			annotationDB,
			annotationJSON,
			suffix,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"",
		"yaml file where schema locations are stored",
	)
	RootCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"",
		"output file",
	)
	RootCmd.Flags().StringVar(
		&annotationDB,
		"db",
		"db",
		"annotation to schemas fields",
	)
	RootCmd.Flags().StringVar(
		&annotationJSON,
		"json",
		"json",
		"annotation to objects fields",
	)
	RootCmd.Flags().StringVar(
		&suffix,
		"suffix",
		"",
		"suffix added to struct names",
	)
}

func Execute() {
	RootCmd.Execute()
}

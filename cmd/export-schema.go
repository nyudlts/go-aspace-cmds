package cmd

import (
	"flag"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"os"
)

var (
	config string
	environment string
	schemaName string
)

func init() {
	schemaCmd.Flags().StringVar(&config,"config", "/etc/go-aspace.yml", "config")
	schemaCmd.Flags().StringVar(&environment, "environment", "", "environment")
	schemaCmd.Flags().StringVar(&schemaName, "schema", "", "schema")
	rootCmd.AddCommand(schemaCmd)
}

var schemaCmd = &cobra.Command {
	Use: "export-schema",
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		client, err := aspace.NewClient(config, environment, 20)
		if err != nil {
			panic(err)
		}

		schema, err := client.GetSchema(schemaName)
		if err != nil {
			panic(err)
		}

		f, err := os.OpenFile(schemaName + ".json", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("error creating output file")
		}
		defer f.Close()

		_,err = f.WriteString(schema)
		if err != nil {
			fmt.Println("error writing to output file")
		}
	},
}

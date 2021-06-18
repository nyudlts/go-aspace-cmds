package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
)

func init() {
	reformatTabsCMD.Flags().StringVar(&location, "location", "", "location of finding aids")
	rootCmd.AddCommand(reformatTabsCMD)
}

var reformatTabsCMD = &cobra.Command{
	Use: "reformat-tabs",
	Run: func(cmd *cobra.Command, args []string) {
		err := reformatTabs(); if err != nil {
			fmt.Println(err)
		}
	},
}

func reformatTabs() error {

	files, err := ioutil.ReadDir(location)
	if err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() != true {
			path := filepath.Join(location, file.Name())
			err := reformatXML(path); if err != nil {
				fmt.Println("Failure: could not reformat ", path)
			}
		}
	}
	return nil
}
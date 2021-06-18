package cmd

import (
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func init() {
	validateDirectoryCMD.Flags().StringVar(&location, "location", "", "Location to validate")
	rootCmd.AddCommand(validateDirectoryCMD)
}

var validateDirectoryCMD = &cobra.Command{
	Use: "validate-directory",
	Run: func(cmd *cobra.Command, args []string) {
		err := validateDirectory(); if err != nil {
			log.Println(err)
		}
	},
}

func validateDirectory() error {
	files, err := ioutil.ReadDir(location)
	if err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() != true {
			path := filepath.Join(location, file.Name())
			log.Println("checking", path)
			filebytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			errorPath := "akkasah-failures"

			err = aspace.ValidateEAD(filebytes); if err != nil {
				fmt.Println(path, "is not valid")
				eadFile := filepath.Join(errorPath, file.Name())
				os.Create(eadFile)
				ioutil.WriteFile(eadFile, filebytes, 0666)
			}
		}
	}

	return nil
}
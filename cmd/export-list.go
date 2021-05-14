package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func init() {
	exportListCmd.Flags().StringVarP(&env, "environment", "e", "dev", "ArchivesSpace environment to be used for export")
	exportListCmd.Flags().StringVarP(&location, "location", "l", ".", "Location to write exported files")
	exportListCmd.Flags().BoolVar(&pretty, "pretty", true, "Pretty format finding aid")
	exportListCmd.Flags().StringVarP(&input, "input", "i", ".", "input file")
	rootCmd.AddCommand(exportListCmd)
}

var exportListCmd = &cobra.Command{
	Use:   "export-list",
	Short: "Export a ead finding aids from a list",
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		client, err := aspace.NewClient(env, 20)
		HandleError(err)

		err = exportList(client)
		if err != nil {
			HandleError(err)
		}
	},
}

func exportList(client *aspace.ASClient) error {
	//check if the export location exists
	if _, err := os.Stat(location); os.IsNotExist(err) {
		return err
	}

	//check if the input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		repositoryId, resourceId, err := aspace.URISplit(scanner.Text())
		if err != nil {
			return err
		}
		fmt.Println("Exporting", scanner.Text())
		resource, err := client.GetResource(repositoryId, resourceId)
		if err != nil {
			return err
		}

		//check that repo location exists, if not create it
		repoPath := filepath.Join(location, shortNames[repositoryId])
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			os.Mkdir(repoPath, 0777)
		}

		//write the ead file to the repo path
		err = getEADFile(repositoryId, resourceId, repoPath, resource.EADID, pretty, client)
		if err != nil {
			log.Println("FAILURE")
		}
	}
	if scanner.Err() != nil {
		return err
	}
	return nil
}

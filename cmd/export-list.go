package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"log"
	"os"
)



func init() {
	exportListCmd.Flags().StringVarP(&env, "environment", "e", "dev", "ArchivesSpace environment to be used for export")
	exportListCmd.Flags().BoolVar(&pretty, "pretty", true, "Pretty format finding aid")
	exportListCmd.Flags().StringVarP(&input, "input", "i", ".", "input file")
	rootCmd.AddCommand(exportListCmd)
}


var exportListCmd = &cobra.Command{
	Use: "export-list",
	Short: "Export a ead finding aids from a list",
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		client, err := aspace.NewClient(env, 20)
		HandleError(err)

		err = exportList(client); if err != nil {
			HandleError(err)
		}
	},
}

func exportList(client *aspace.ASClient) error {
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
		err = getEADFile(repositoryId, resourceId, location, resource.EADID, pretty, client); if err != nil {
			log.Println("FAILURE")
		}
	}
	if scanner.Err() != nil {
		return err
	}
	return nil
}

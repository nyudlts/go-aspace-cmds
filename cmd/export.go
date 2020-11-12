package cmd

import (
	"bufio"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var env string
var repositoryId int
var resourceId int

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export an EAD from ArchivesSpace with go-aspace",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := aspace.NewClient(env, 20)
		HandleError(err)
		err = exportEAD(client)
		HandleError(err)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&env, "environment", "e", "dev", "ArchivesSpace environment to be used for export")
	exportCmd.Flags().IntVar(&repositoryId, "repository", 2, "Repository to be used for export")
	exportCmd.Flags().IntVar(&resourceId, "resource", 1, "Resource to be exported")
	rootCmd.AddCommand(exportCmd)
}

func exportEAD(client *aspace.ASClient) error {
	log.Println("go-aspace lib", aspace.LibraryVersion)

	//request the resource
	resource, err := client.GetResource(repositoryId, resourceId)
	if err != nil {
		return err
	}

	//create filename from resource ids
	outputTitle := resource.ID0
	for _, id := range []string{resource.ID1, resource.ID2, resource.ID3} {
		if id != "" {
			outputTitle = outputTitle + "-" + id
		}
	}
	log.Println("Exporting", outputTitle)
	outputTitle = outputTitle + ".xml"

	//request the ead of the resource
	ead, err := client.SerializeEAD(repositoryId, resourceId, true, false, false, false, false)
	if err != nil {
		return err
	}

	//Validate the ead
	err = aspace.ValidateEAD(ead)
	if err != nil {
		log.Println("WARNING: Exported EAD file did not pass validation")
	}

	//create the output file
	output, err := os.Create(outputTitle)
	if err != nil {
		return err
	}
	defer output.Close()

	//write ead []bytes to file
	writer := bufio.NewWriter(output)
	_, err = writer.Write(ead)
	if err != nil {
		return err
	}
	writer.Flush()

	//done
	log.Println("Export complete")
	return nil
}

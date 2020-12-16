package cmd

import (
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

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
			outputTitle = outputTitle + "_" + id
		}
	}

	outputTitle = strings.ToLower(outputTitle)
	log.Println("Exporting", outputTitle)
	outputTitle = outputTitle + ".xml"

	err = getEADFile(repositoryId, resourceId, ".", resource.EADID, client)
	if err != nil {
		panic(err)
	}

	//done
	log.Printf("Export of %s complete\n", resource.EADID)
	return nil
}

package cmd

import (
	"bufio"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
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
	fmt.Println("go-aspace lib", aspace.LibraryVersion)
	fmt.Println(fmt.Sprintf("Exporting /repostories/%d/resources/%d from environment %s", repositoryId, resourceId, env))

	ead, err := client.SerializeEAD(repositoryId, resourceId, true, false, false, false, false)
	if err != nil {
		return err
	}

	err = aspace.ValidateEAD(ead)
	if err != nil {
		return err
	}

	output, err := os.Create(fmt.Sprintf("%d_%d.xml", repositoryId, resourceId))
	if err != nil {
		return err
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	_, err = writer.Write(ead)
	if err != nil {
		return err
	}
	writer.Flush()

	return nil
}

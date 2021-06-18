package cmd

import (
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func init() {
	exportResourceCmd.Flags().StringVarP(&env, "environment", "e", "dev", "ArchivesSpace environment to be used for export")
	exportResourceCmd.Flags().IntVar(&repositoryId, "repository", 0, "Repository to be used for export")
	exportResourceCmd.Flags().IntVar(&resourceId, "resource", 0, "Resource to be exported")
	exportResourceCmd.Flags().StringVarP(&location, "location", "l", ".", "location to export finding aids")
	exportResourceCmd.Flags().BoolVar(&pretty, "pretty", false, "Pretty format finding aid")
	rootCmd.AddCommand(exportResourceCmd)
}

var exportResourceCmd = &cobra.Command{
	Use:   "export-resource",
	Short: "export an EAD from ArchivesSpace with go-aspace",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		client, err = aspace.NewClient("/etc/go-aspace.yml", env, 20)
		if err != nil {
			panic(err)
		}

		err = ExportResource()
		if err != nil {
			panic(err)
		}
	},
}

func ExportResource() error {
	repository, err := client.GetRepository(repositoryId)
	if err != nil {
		return fmt.Errorf("Repistory ID %d does not exist", repository)
	}

	slug := repository.Slug

	outputDir := filepath.Join(location, slug)

	err = os.MkdirAll(outputDir, 0775)
	if err != nil {
		return fmt.Errorf("Could not create an export directory for repository %s", repository)
	}

	err = exportEAD(resourceId, outputDir, "failures")
	if err != nil {
		return err
	}

	return nil
}

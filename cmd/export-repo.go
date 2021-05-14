package cmd

import (
	"flag"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func init() {
	exportRepoCmd.Flags().StringVarP(&env, "environment", "e", "", "ArchivesSpace environment to be used for export")
	exportRepoCmd.Flags().IntVarP(&repositoryId, "repository", "r", 2, "Repository to be used for export")
	exportRepoCmd.Flags().StringVarP(&location, "location", "l", ".", "location to export finding aids")
	exportRepoCmd.Flags().BoolVar(&pretty, "pretty", false, "Pretty format finding aid")
	rootCmd.AddCommand(exportRepoCmd)
}

var exportRepoCmd = &cobra.Command{
	Use:   "export-repo",
	Short: "Export all published finding aids",
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		var err error
		client, err = aspace.NewClient(env, 20)
		HandleError(err)

		fmt.Printf("go-aspace version %s\n", aspace.LibraryVersion)

		err = exportRepo(repositoryId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func exportRepo(repoId int) error {

	repository, err := client.GetRepository(repositoryId)
	if err != nil {
		return fmt.Errorf("Repistory ID %d does not exist", repoId)
	}

	slug := repository.Slug

	outputDir := filepath.Join(location, slug)

	err = os.MkdirAll(outputDir, 0775)
	if err != nil {
		return fmt.Errorf("Could not create an export directory for repository %s", repoId)
	}

	resourceIds, err := client.GetResourceIDs(repoId)
	if err != nil {
		return fmt.Errorf("Could not get resource List for repository %d", repoId)
	}

	if len(resourceIds) <= 0 {
		return fmt.Errorf("Repository '%s' does not contain any resources to export", slug)
	}

	for _, resourceId := range resourceIds {
		err = exportEAD(resourceId, outputDir)
		if err != nil {
			log.Println(err.Error())
		}

	}

	return nil
}

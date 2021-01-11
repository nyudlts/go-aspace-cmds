package cmd

import (
	"bufio"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var repoID int

func init() {
	exportRepoCmd.Flags().IntVarP(&repoID, "repository", "r", 2,"The Id of the Repository to be Exported")
	rootCmd.AddCommand(exportRepoCmd)
}

var exportRepoCmd = &cobra.Command{
	Use:   "export-repo",
	Short: "Export all published finding aids",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		client, err := aspace.NewClient("fade", 20)
		HandleError(err)

		err = exportRepo(repoID, client)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("complete, run time:", time.Since(start))

	},
}

func exportRepo(repoId int, client *aspace.ASClient) error {
	repository, err := client.GetRepository(repoID)
	if err != nil {
		return err
	}

	err = os.Mkdir("exports", 0777)
	if err != nil {
		return err
	}

	logFile, err := os.Create(repository.Slug + ".txt")
	if err != nil {
		return err
	}
	defer logFile.Close()
	writer := bufio.NewWriter(logFile)

	resourceIds, err := client.GetResourceIDs(repoId)
	if err != nil {
		return fmt.Errorf("Could not get resource List for repository %d", repoId)
	}

	for _, resourceId := range resourceIds {
		resource, err := client.GetResource(repoId, resourceId)
		if err != nil {
			fmt.Printf(fmt.Sprintf("Could not get resource %d from repo %d, skipping", resourceId, repoId))
		}

		if resource.Publish == true {
			log.Println("attempting ", resource.EADID, resource.URI)
			err = getEADFile(repoId, resourceId, "exports", resource.EADID, client)
			if err != nil {
				log.Println(err)
				writer.WriteString(fmt.Sprintf("%s\t%s\t%v\n", resource.URI, "failure", err))
				writer.Flush()
			} else {
				log.Println("success")
				writer.WriteString(fmt.Sprintf("%s\t%s\n", resource.URI, "success"))
				writer.Flush()
			}
		}

	}

	return nil
}

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
	allCmd.Flags().StringVarP(&env, "environment", "e", "dev", "ArchivesSpace environment to be used for export")
	allCmd.Flags().IntVarP(&repositoryId, "repository", "r",2, "Repository to be used for export")
	rootCmd.AddCommand(allCmd)
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Export all published finding aids",
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		client, err := aspace.NewClient(env, 20)
		HandleError(err)

		err = exportRepo(repositoryId, client)
		if err != nil {
			fmt.Println(err)
		}

	},
}

func exportRepo(repoId int, client *aspace.ASClient) error {
	repoString := fmt.Sprintf("%d", repositoryId)
	out, err := os.Create(fmt.Sprintf("errors-repository-%s.txt", repoString))
	if err != nil {
		return err
	}
	defer out.Close()
	writer := bufio.NewWriter(out)

	resourceIds, err := client.GetResourceIDs(repoId)
	if err != nil {
		return fmt.Errorf("Could not get resource List for repository %d", repoId)
	}

	err = os.Mkdir(repoString, 0775); if err != nil {
		return fmt.Errorf("Could not create an export directory for repository %d", repoId)
	}

	for _, resourceId := range resourceIds {
		resource, err := client.GetResource(repoId, resourceId)
		if err != nil {
			fmt.Printf(fmt.Sprintf("Could not get resource %d from repo %d, skipping", resourceId, repoId))
		}

		if resource.Publish == true {
			log.Println("attempting", resource.EADID, resource.URI)
			err = getEADFile(repoId, resourceId, repoString, resource.EADID, client)
			if err != nil {
				fmt.Println(err)
				writer.WriteString(fmt.Sprintf("%s\t%v\n", resource.URI, err))
				writer.Flush()
			}
		}

	}

	return nil
}

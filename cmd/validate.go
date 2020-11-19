package cmd

import (
	"bufio"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var validateCmd = &cobra.Command{
	Use: "validate",
	Short: "Validate ead exports for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := aspace.NewClient(env, timeout)
		if err != nil {
			HandleError(err)
		}
		validateRepository(client)
	},
}

func init() {
	validateCmd.Flags().StringVarP(&env, "environment", "e", "dev", "ArchivesSpace environment to be used for export")
	validateCmd.Flags().IntVar(&repositoryId, "repository", 2, "Repository to be used for export")
	validateCmd.Flags().IntVar(&timeout, "timeout", 20, "server timeout")
	rootCmd.AddCommand(validateCmd)
}

func validateRepository(client *aspace.ASClient) {
	fmt.Println("go-aspace", aspace.LibraryVersion)
	fmt.Printf("Validating resources in repository %i\n", repositoryId)
	outFile, err := os.Create("validation-failures.txt")
	HandleError(err)
	writer := bufio.NewWriter(outFile)

	resourceIds, err := client.GetResourceIDs(repositoryId)
	HandleError(err)

	for i, resourceId := range resourceIds {
		resource, err := client.GetResource(repositoryId, resourceId)
		HandleError(err)
		fmt.Print(i, " ", resource.Title)
		ead, err := client.SerializeEAD(repositoryId, resourceId, true, false, false, false, false)
		err = aspace.ValidateEAD(ead); if err != nil {
			fmt.Print( "\t *Failed Validation*\n")
			writer.WriteString(resource.URI +"\n")
			writer.Flush()
			outEAD, innerErr := os.Create(filepath.Join("failures", resource.EADID + ".xml"))
			HandleError(innerErr)
			eadWriter := bufio.NewWriter(outEAD)
			eadWriter.Write(ead)
			eadWriter.Flush()
		} else {
			fmt.Print("\t *Passed Validation*\n")
		}
	}

}
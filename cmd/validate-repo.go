package cmd

import (
	"bufio"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var validationType string
var published bool

var validateCmd = &cobra.Command{
	Use:   "validate-repo",
	Short: "Validate ead exports for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		client, err = aspace.NewClient("/etc/go-aspace.yml", env, 20)
		if err != nil {
			HandleError(err)
		}

		validateRepository()


	},
}

func init() {
	validateCmd.Flags().StringVarP(&env, "environment", "e", "dev", "ArchivesSpace environment to be used for export")
	validateCmd.Flags().IntVar(&repositoryId, "repository", 2, "Repository Id to be used for validated")
	validateCmd.Flags().BoolVar(&published, "published", false, "Export only published resources")
	validateCmd.Flags().BoolVar(&pretty, "pretty", false, "Pretty format finding aid")
	rootCmd.AddCommand(validateCmd)
}

func validateRepository() {
	fmt.Println("go-aspace", aspace.LibraryVersion)
	fmt.Printf("Validating resources in repository %d\n", repositoryId)
	outFile, err := os.Create("validation-failures.txt")
	HandleError(err)
	writer := bufio.NewWriter(outFile)

	resourceIds, err := client.GetResourceIDs(repositoryId)
	HandleError(err)

	for i, resId := range resourceIds {
		resource, err := client.GetResource(repositoryId, resId)
		HandleError(err)

		if published == resource.Publish {
			fmt.Print(i, " ", resource.Title)
			ead, err := client.SerializeEAD(repositoryId, resId, true, false, false, false, false)
			err = aspace.ValidateEAD(ead)
			if err != nil {
				fmt.Print("\t *Failed Validation*\n")
				writer.WriteString(resource.URI + "\n")
				writer.Flush()
				outEAD, innerErr := os.Create(filepath.Join("failures", resource.EADID+".xml"))
				HandleError(innerErr)
				eadWriter := bufio.NewWriter(outEAD)
				eadWriter.Write(ead)
				eadWriter.Flush()
			} else {
				fmt.Print("\t *Passed Validation*\n")
			}
		}
	}
}

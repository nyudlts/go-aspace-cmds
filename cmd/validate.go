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
	Use:   "validate",
	Short: "Validate ead exports for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := aspace.NewClient(env, timeout)
		if err != nil {
			HandleError(err)
		}
		if validationType == "repository" {
			validateRepository(client)
		} else if validationType == "resource" {
			validateResource(client)
		} else {
			panic(fmt.Errorf("type argument must be either 'repository' or 'resource'"))
		}

	},
}

func init() {
	validateCmd.Flags().StringVarP(&env, "environment", "e", "dev", "ArchivesSpace environment to be used for export")
	validateCmd.Flags().IntVar(&repositoryId, "repository", 2, "Repository Id to be used for validated")
	validateCmd.Flags().IntVar(&resourceId, "resource", 1, "Resource Id to be validated")
	validateCmd.Flags().IntVar(&timeout, "timeout", 20, "server timeout")
	validateCmd.Flags().StringVarP(&validationType, "type", "t", "", "type of validation to perform")
	validateCmd.Flags().BoolVar(&published, "published", true, "Export only published resources")
	exportCmd.Flags().BoolVar(&pretty, "pretty", false, "Pretty format finding aid")
	rootCmd.AddCommand(validateCmd)
}

func validateRepository(client *aspace.ASClient) {
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

func validateResource(client *aspace.ASClient) {
	resource, err := client.GetResource(repositoryId, resourceId)
	HandleError(err)
	ead, err := client.SerializeEAD(repositoryId, resourceId, true, false, false, false, false)
	HandleError(err)
	err = aspace.ValidateEAD(ead)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resource.Title, " passed validation")
	}
}

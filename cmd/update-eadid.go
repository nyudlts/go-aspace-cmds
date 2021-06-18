package cmd

import (
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	updateEadidCmd.Flags().StringVarP(&env, "environment", "e", "", "ArchivesSpace environment to be used for export")
	updateEadidCmd.Flags().IntVarP(&repositoryId, "repository", "r", 0, "Repository to be used for export")
	rootCmd.AddCommand(updateEadidCmd)
}

var updateEadidCmd = &cobra.Command{
	Use: "update-eadid",
	Short: "Update eadid to publisher 2.0 spec",
	Run: func(cmd *cobra.Command, args []string) {
		err := updateEadid(); if err != nil {
			fmt.Println (err.Error())
		}
	},
}

func updateEadid() error {

	var err error
	client, err = aspace.NewClient("/etc/go-aspace.yml", env, 20)
	if err != nil {
		return err
	}

	resourceIds, err := client.GetResourceIDs(repositoryId)
	if err != nil {
		return err
	}

	for _, resourceId := range resourceIds {
		resource, err := client.GetResource(repositoryId, resourceId)
		if err != nil {
			return err
		}

		eadid := resource.EADID
		eids := resource.ID0
		if resource.ID1 != "" {
			eids = eids + "_" + resource.ID1
		}
		if resource.ID2 != "" {
			eids = eids + "_" + resource.ID2
		}
		if resource.ID3 != "" {
			eids = eids + "_" + resource.ID3
		}

		eids = strings.ToLower(eids)

		if eadid != eids {
			 fmt.Println("updating", eadid, "to", eids)
			 resource.EADID = eids
			 msg, err := client.UpdateResource(repositoryId, resourceId, resource)
			 if err != nil {
			 	return err
			 }
			 fmt.Println(msg)
		}
	}

	return nil
}


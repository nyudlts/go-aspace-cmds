package cmd

import (
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

var count int

var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Generate a sample set of ead finding aids",
	Run: func(cmd *cobra.Command, args []string) {
		err := sample()
		HandleError(err)
	},
}

func init() {
	sampleCmd.Flags().IntVar(&count, "count", 1, "number of finding aids to export")
	rootCmd.AddCommand(sampleCmd)
}

func sample() error {
	validEadCount := 0
	//initialize stuff
	var s1 = rand.NewSource(time.Now().UnixNano())
	var r1 = rand.New(s1)

	//get a client
	client, err := aspace.NewClient("fade", 20)
	if err != nil {
		return err
	}

	//get a list of repos
	repoIds, err := client.GetRepositories()
	if err != nil {
		return err
	}

	//make a map of resource ids for each repository
	repositoryResourceIds := map[int][]int{}
	for _, repoId := range repoIds {
		rIds, err := client.GetResourceIDs(repoId)
		if err != nil {
			return err
		}
		repositoryResourceIds[repoId] = rIds
	}

	//generate random finding aids
	for validEadCount < count {
		repoId := repoIds[r1.Intn(len(repoIds)-3)]
		resourceIds := repositoryResourceIds[repoId]
		resourceId := resourceIds[r1.Intn(len(resourceIds))]
		fmt.Println("attempting to serialize", repoId, resourceId)
		if GetEADFile(repoId, resourceId, "", "TEMP", client) == nil {
			validEadCount = validEadCount + 1
		}
	}

	return nil
}

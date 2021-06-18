package cmd

import (
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"log"
)

var useStatements = map[string]int{}

func contains(s string) bool {
	for k,_ := range useStatements {
		if k == s {
			return true
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(UseStatementCmd)
}

var UseStatementCmd = &cobra.Command{
	Use: "use-statements",
	Run: func(cmd *cobra.Command, args []string) {


		log.Println("INFO Analyzing Use Statements")
		client, err := aspace.NewClient("/etc/go-aspace.yml", "fade", 20)
		if err != nil {
			panic(err)
		}

		log.Println("INFO go-aspace", aspace.LibraryVersion)

		for _, repoId := range []int{2,3,6} {
			log.Println("INFO Getting DOs for repository", repoId)

			doIds, err := client.GetDigitalObjectIDs(repoId)
			if err != nil {
				panic(err)
			}
			for i,doId := range doIds {
				if i % 1000 == 0 {
					log.Println("INFO", i, "of", len(doIds), "records processed")
				}

				do, err := client.GetDigitalObject(repoId, doId)
				if err != nil {
					log.Println("ERROR could not get do", doId, "in repo", repoId)
				}

				if len(do.FileVersions) > 1 {
					log.Println("INFO", "MULTIPLE USES: " + do.URI)
				} else if len(do.FileVersions) == 1 {
					if contains(do.FileVersions[0].UseStatement) == true {
						useStatements[do.FileVersions[0].UseStatement] = useStatements[do.FileVersions[0].UseStatement]  + 1
					} else {
						useStatements[do.FileVersions[0].UseStatement] = 1
					}
				}
			}

			fmt.Println(repoId, "\n-----")

			for k,v := range useStatements {
				fmt.Println(k,v)
			}

			useStatements = map[string]int{}
		}
	},
}

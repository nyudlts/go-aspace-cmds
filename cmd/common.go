package cmd

import (
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
)

var (
	client *aspace.ASClient
	env string
	repositoryId int
	resourceId int
	timeout int
)

var rootCmd = &cobra.Command{
	Use:   "go-aspace",
	Short: "A tool to run go-aspace scripts",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		HandleError(err)
	}
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

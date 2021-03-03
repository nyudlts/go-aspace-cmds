package cmd

import (
	"bufio"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	client       *aspace.ASClient
	env          string
	repositoryId int
	resourceId   int
	timeout      int
	location	string
	pretty		bool
	input		string
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

func getEADFile(repoId int, resourceId int, loc string, eadid string, pretty bool, client *aspace.ASClient) error {

	ead, err := client.GetEADAsByteArray(repoId, resourceId)
	if err != nil {
		return err
	}

	if len(ead) <= 0 {
		return fmt.Errorf("Returned a zero length array")
	}

	err = aspace.ValidateEAD(ead)
	if err != nil {
		return err
	}

	file := filepath.Join(loc, fmt.Sprintf("%s.xml", eadid))
	eadFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer eadFile.Close()

	writer := bufio.NewWriter(eadFile)
	_, err = writer.Write(ead)
	if err != nil {
		os.Remove(file)
		return err
	}

	if pretty == true {
		out, err := exec.Command("xmllint", "--format", file).Output()
		if err != nil {
			panic(err)
		}

		eadFile.Close()
		f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		writer = bufio.NewWriter(f)
		writer.Write(out)
		writer.Flush()
	}
	return nil
}
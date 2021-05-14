package cmd

import (
	"fmt"
	"github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"io/ioutil"
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
	location     string
	pretty       bool
	input        string
)

var shortNames = map[int]string{
	2: "tamwag",
	3: "fales",
	6: "archives",
}
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

func exportEAD(resId int, outputDir string) error {
	resource, err := client.GetResource(repositoryId, resId)
	if err != nil {
		msg := fmt.Errorf("WARNING Could not get resource %d from repo %d, skipping", resourceId, repositoryId)
		return msg
	}

	if resource.Publish == true {

		eadFilename := resource.EADID
		if resource.EADID == "" {
			eadFilename = fmt.Sprintf("%d_%d", repositoryId, resId)
			log.Println(fmt.Sprintf("WARNING %s does not have an eadid, substituting %s as the filename.", resource.URI, eadFilename))
		}

		log.Println("INFO exporting", eadFilename+".xml", resource.URI)

		err = getEADFile(repositoryId, resId, outputDir, eadFilename, pretty, client)
		if err != nil {
			log.Println("ERROR", err.Error())
		} else {
			log.Println("INFO", eadFilename+".xml exported")
		}
	}

	return nil
}

func getEADFile(repoId int, resourceId int, loc string, eadid string, pretty bool, client *aspace.ASClient) error {

	ead, err := client.GetEADAsByteArray(repoId, resourceId)
	if err != nil {
		return fmt.Errorf("ArchiveSpace did not return an EAD file for %s", eadid)
	}

	if len(ead) <= 0 {
		return fmt.Errorf("ArchiveSpace returned a zero length array")
	}

	err = aspace.ValidateEAD(ead)
	if err != nil {
		return fmt.Errorf("EAD validation failed on %s", eadid)
	}

	exportFile := fmt.Sprintf("%s.xml", eadid)
	file := filepath.Join(loc, exportFile)

	err = ioutil.WriteFile(file, ead, 0644)
	if err != nil {
		return fmt.Errorf("Could not write to file %s", file)
	}

	if pretty == true {
		err = reformatXML(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func reformatXML(path string) error {

	reformattedBytes, err := exec.Command("xmllint", "--format", path).Output()
	if err != nil {
		return fmt.Errorf("could not reformat %s", path)
	}

	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("could not delete %s", path)
	}

	err = ioutil.WriteFile(path, reformattedBytes, 0644)
	if err != nil {
		return fmt.Errorf("could not write reformated bytes to %s", path)
	}

	return nil
}

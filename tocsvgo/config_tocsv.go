package tocsvgo

import (
	"fmt"

	"github.com/parmaanu/goutils/errorutils"
	"github.com/parmaanu/goutils/filesystem"
	"github.com/parmaanu/goutils/fileutils"

	"github.com/goccy/go-yaml"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

// TocsvConfig stores the config related to the Tocsv application
type TocsvConfig struct {
	AnchorFiles      []string `yaml:"AnchorFiles"`
	PrintTagInOutput bool     `yaml:"PrintTagInOutput"`
	LogDirectory     string   `yaml:"LogDirectory"`
}

// NewToCsvConfig return an instance of TocsvConfig struct
func NewToCsvConfig(configFile string) *TocsvConfig {
	absfname, _ := tilde.Expand(configFile)
	if !fileutils.FileExist(absfname) {
		fmt.Println("ERROR:", absfname, "does not exist")
		return nil
	}
	reader, err := filesystem.Open(absfname)
	errorutils.PanicOnErr(err)

	decoder := yaml.NewDecoder(reader)

	// don't check for error while decoding the config here as we don't know the alias files
	tocsvConfig := &TocsvConfig{}
	decoder.Decode(tocsvConfig)
	return tocsvConfig
}

package filterlogs_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"tocsv/filterlogs"

	"github.com/parmaanu/goutils/filesystem"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var gStdoutReader, gStdoutWriter, gOriginalStdout *os.File
var gConfigFile, gFname string
var gAnchorFiles []string

func init() {
	mfs := filesystem.NewMockFileSystem()
	filesystem.SetFileSystem(mfs)

	gFname = "test.log"
	mfs.SetFileData(gFname, []string{"2020-06-02 14:33:56.531063 ORDER NEW price: 123.123, quantity: 1000, securityId: 999, side: BUY, bid: 124.0, ask: 125.0"})

	var mainConfigStr = `
Apps:
	- *Orders
	`
	gConfigFile = "full_filterlogs.yaml"
	mfs.SetFileData(gConfigFile, strings.Split(mainConfigStr, "\n"))

	// Note, since anchor file is used as `yaml.ReferenceFiles(absAnchorFiles...)`, its data cannot be set in MockFileSystem
	gAnchorFiles = []string{"orders.yaml"}
	for _, f := range gAnchorFiles {
		mfs.SetFileData(f, []string{})
	}
}

func captureStdout() {
	gStdoutReader, gStdoutWriter, _ = os.Pipe()
	gOriginalStdout = os.Stdout // save original stdout
	os.Stdout = gStdoutWriter
}

func getCapturedStdout() string {
	gStdoutWriter.Close()
	out, _ := ioutil.ReadAll(gStdoutReader)
	os.Stdout = gOriginalStdout // restore original stdout
	return string(out)
}

func TestFilterLogsWithSingleApp(t *testing.T) {
	interactiveMode := false

	expectedFilteredData := map[string]*filterlogs.FilteredData{
		"TimeStampKey":  {"2020-06-02 14:33:56.531063"},
		"SecurityIdKey": {"999"},
		"PriceKey":      {"123.123"},
		"QuantityKey":   {"1000"},
		"SideKey":       {"BUY"},
		"BidKey":        {"124.0"},
		"AskKey":        {"125.0"},
	}

	callbackCalled := false
	app := filterlogs.NewApp([]string{gFname}, gConfigFile, gAnchorFiles, interactiveMode)
	callback := func(config *filterlogs.ClientConfigType, filteredData map[string]*filterlogs.FilteredData) {
		assert.Equal(t, "Orders", config.AppName)

		if diff := cmp.Diff(expectedFilteredData, filteredData); diff != "" {
			t.Errorf("FilteredData map not equal:\n%s", diff)
		}
		callbackCalled = true
	}
	app.Run(callback)
	assert.True(t, callbackCalled, "Assigned callback is not called. Please check the config.")
}

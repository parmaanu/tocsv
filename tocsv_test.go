package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/parmaanu/goutils/filesystem"
	"github.com/parmaanu/goutils/fileutils"

	"github.com/stretchr/testify/assert"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

var gMfs *filesystem.MockFileSystem
var gStdoutReader, gStdoutWriter, gOriginalStdout *os.File

func init() {
	gMfs = filesystem.NewMockFileSystem()
	filesystem.SetFileSystem(gMfs)
	os.Unsetenv("SHOWCSV_RENDER_TUI")
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

func TestToCsvWithSingleApp_and_PrintTagInOutput_to_false(t *testing.T) {
	captureStdout()

	fname := "test.log"
	gMfs.SetFileData(fname, []string{"2020-06-02 14:33:56.531063 ORDER NEW price: 123.123, quantity: 1000, securityId: 999, side: BUY, bid: 124.0, ask: 125.0"})
	configFile := "tocsv.yaml"
	anchorFiles := []string{"orders.yaml", "columns.yaml"}
	printOnStdout := true
	interactiveMode := false
	// Note, PrintTagInOutput is not given in the config file so default value would be false
	gMfs.SetFileData(configFile, []string{`
Apps:
	- *Orders
	`})

	for _, f := range anchorFiles {
		gMfs.SetFileData(f, []string{})
	}

	tocsv := NewTocsv([]string{fname}, configFile, anchorFiles, printOnStdout, interactiveMode)
	assert.NotNil(t, tocsv, "not able to create tocsv instance")
	tocsv.Run()

	output := getCapturedStdout()
	expectedOutput := `timestamp,securityId,price,quantity,side,bid,ask
2020-06-02 14:33:56.531063,999,123.123,1000,BUY,124.0,125.0
`
	assert.Equal(t, expectedOutput, output)
}

func Test_tocsv_prints_tag_column_in_output_if_PrintTagInOutput_is_set_to_true(t *testing.T) {
	captureStdout()

	fname := "test.log"
	gMfs.SetFileData(fname, []string{"2020-06-02 14:33:56.531063 ORDER NEW price: 123.123, quantity: 1000, securityId: 999, side: BUY, bid: 124.0, ask: 125.0"})
	configFile := "tocsv.yaml"
	anchorFiles := []string{"orders.yaml", "columns.yaml"}
	printOnStdout := true
	interactiveMode := false
	gMfs.SetFileData(configFile, []string{`
PrintTagInOutput: true
Apps:
	- *Orders
	`})

	for _, f := range anchorFiles {
		gMfs.SetFileData(f, []string{})
	}

	tocsv := NewTocsv([]string{fname}, configFile, anchorFiles, printOnStdout, interactiveMode)
	assert.NotNil(t, tocsv, "not able to create tocsv instance")
	tocsv.Run()

	output := getCapturedStdout()
	expectedOutput := `__tag__,timestamp,securityId,price,quantity,side,bid,ask
NEW,2020-06-02 14:33:56.531063,999,123.123,1000,BUY,124.0,125.0
`
	assert.Equal(t, expectedOutput, output)
}

func Test_tocsv_when_LogDirectory_is_given_in_config(t *testing.T) {
	fname := "test.log"
	gMfs.SetFileData(fname, []string{"2020-06-02 14:33:56.531063 ORDER NEW price: 123.123, quantity: 1000, securityId: 999, side: BUY, bid: 124.0, ask: 125.0"})
	configFile := "tocsv.yaml"
	anchorFiles := []string{"orders.yaml", "columns.yaml"}
	printOnStdout := false
	interactiveMode := false
	gMfs.SetFileData(configFile, []string{`
PrintTagInOutput: true
LogDirectory: ~/logs/
Apps:
	- *Orders
	`})

	for _, f := range anchorFiles {
		gMfs.SetFileData(f, []string{})
	}

	tocsv := NewTocsv([]string{fname}, configFile, anchorFiles, printOnStdout, interactiveMode)
	assert.NotNil(t, tocsv, "not able to create tocsv instance")
	tocsv.Run()

	dt := time.Now().Format("20060102") // YYYYMMDD format
	outputFileName, err := tilde.Expand(fmt.Sprintf("~/logs/Orders.%s.csv", dt))
	assert.NoError(t, err, "error found while expanding full path")
	assert.FileExists(t, outputFileName, outputFileName+" file does not exist")

	output := string(fileutils.ReadFullFileAsBytes(outputFileName))
	expectedOutput := `__tag__,timestamp,securityId,price,quantity,side,bid,ask
NEW,2020-06-02 14:33:56.531063,999,123.123,1000,BUY,124.0,125.0
`
	assert.Equal(t, expectedOutput, output)
	err = os.Remove(outputFileName)
	assert.NoError(t, err, "error while removing the output file")
}

func Test_tocsv_when_LogDirectory_is_not_given_in_config(t *testing.T) {
	fname := "test.log"
	gMfs.SetFileData(fname, []string{"2020-06-02 14:33:56.531063 ORDER NEW price: 123.123, quantity: 1000, securityId: 999, side: BUY, bid: 124.0, ask: 125.0"})
	configFile := "tocsv.yaml"
	anchorFiles := []string{"orders.yaml", "columns.yaml"}
	printOnStdout := false
	interactiveMode := false
	gMfs.SetFileData(configFile, []string{`
PrintTagInOutput: true
Apps:
	- *Orders
	`})

	for _, f := range anchorFiles {
		gMfs.SetFileData(f, []string{})
	}

	tocsv := NewTocsv([]string{fname}, configFile, anchorFiles, printOnStdout, interactiveMode)
	assert.NotNil(t, tocsv, "not able to create tocsv instance")
	tocsv.Run()

	dt := time.Now().Format("20060102") // YYYYMMDD format
	outputFileName, err := tilde.Expand(fmt.Sprintf("./Orders.%s.csv", dt))
	assert.NoError(t, err, "error found while expanding full path")
	assert.FileExists(t, outputFileName, "file does not exist")

	output := string(fileutils.ReadFullFileAsBytes(outputFileName))
	expectedOutput := `__tag__,timestamp,securityId,price,quantity,side,bid,ask
NEW,2020-06-02 14:33:56.531063,999,123.123,1000,BUY,124.0,125.0
`
	assert.Equal(t, expectedOutput, output)
	err = os.Remove(outputFileName)
	assert.NoError(t, err, "error while removing the output file")
}

func Test_AnchorFiles_are_read_through_config_file(t *testing.T) {
	captureStdout()

	fname := "test.log"
	gMfs.SetFileData(fname, []string{"2020-06-02 14:33:56.531063 ORDER NEW price: 123.123, quantity: 1000, securityId: 999, side: BUY, bid: 124.0, ask: 125.0"})
	configFile := "tocsv.yaml"
	// Note, indention is very important in yaml. So, if AnchorFiles and Apps should start with 0 indent
	gMfs.SetFileData(configFile, []string{`
AnchorFiles:
	- orders.yaml
	- columns.yaml

Apps:
	- *Orders
	`})

	// Note, anchorFiles as command line arguments is an empty slice. It is read from the config file
	anchorFiles := []string{}
	printOnStdout := true
	interactiveMode := false

	gMfs.SetFileData("orders.yaml", []string{})
	gMfs.SetFileData("columns.yaml", []string{})

	tocsv := NewTocsv([]string{fname}, configFile, anchorFiles, printOnStdout, interactiveMode)
	assert.NotNil(t, tocsv, "not able to create tocsv instance")
	tocsv.Run()

	output := getCapturedStdout()
	expectedOutput := `timestamp,securityId,price,quantity,side,bid,ask
2020-06-02 14:33:56.531063,999,123.123,1000,BUY,124.0,125.0
`
	assert.Equal(t, expectedOutput, output)
}

func Test_AnchorFiles_given_through_commandline_overrides_AnchorFiles_given_in_config_file(t *testing.T) {
	captureStdout()

	fname := "test.log"
	gMfs.SetFileData(fname, []string{"2020-07-12 01:54:23.124127 POSITION securityId: 154, netPosition: -1230, startOfDayPosition: 1240, dayTradedVolume: 2470"})
	configFile := "tocsv.yaml"
	// Note, indention is very important in yaml. So, if AnchorFiles and Apps should start with 0 indent
	gMfs.SetFileData(configFile, []string{`
AnchorFiles:
- orders.yaml

Apps:
- *Position
`})
	anchorFiles := []string{"position.yaml", "columns.yaml"}
	printOnStdout := true
	interactiveMode := false
	for _, f := range anchorFiles {
		gMfs.SetFileData(f, []string{})
	}

	tocsv := NewTocsv([]string{fname}, configFile, anchorFiles, printOnStdout, interactiveMode)
	assert.NotNil(t, tocsv, "not able to create tocsv instance")
	tocsv.Run()

	output := getCapturedStdout()
	expectedOutput := `timestamp,securityId,netPosition,sodPosition,dayTradedVolume
2020-07-12 01:54:23.124127,154,-1230,1240,2470` + "\n"
	assert.Equal(t, expectedOutput, output)
}

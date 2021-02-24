package filterlogs

import (
	"fmt"
	"strings"
	"tocsv/logparser"

	"github.com/parmaanu/goutils/algoutils"
	"github.com/parmaanu/goutils/fileutils"
)

// // TODO,
// // - remove the tags or do not print the tags which are not present in
// //	 the logfile. This will make number of output columns smaller
// // - Create a config generate from logs
// // - Create a tui to select output flags or columns
// // - Use goevaluate to add a calculated columns
// // - Show log patterns and log line examples quickly using an interactive menu [DONE]
// // - Tail logs for a pattern
// // - Support `After` in log line pattern config to pick up other patterns when there are multiple matching patterns
// in single line
//
// // TODO, golang
// // - create template for creating golang apps
//

const (
	gStartOfLine = "^"
	gEndOfLine   = "$"
)

// FilteredData store the config for each extracted value
type FilteredData struct {
	Text string
}

// ClientCallbackType is the type of the function that is called when we get some filtered data
type ClientCallbackType func(*ClientConfigType, map[string]*FilteredData)

// App is a struct which converts a logfile into a csv
type App struct {
	inputFile   string
	configFile  string
	anchorFiles []string
	interactive bool

	// header    []string
	clientValuesMap map[string]*FilteredData
	clientCallback  ClientCallbackType
}

// NewApp returns an instance of to csv app
func NewApp(inputFiles []string, configFile string, anchorFiles []string, interactiveMode bool) *App {
	return &App{
		// TODO, add a support of multiple input log files later on - need to enhance logparser to read multiple
		// logfiles in sorted order
		inputFile:   inputFiles[0],
		configFile:  configFile,
		anchorFiles: anchorFiles,
		interactive: interactiveMode,

		clientValuesMap: make(map[string]*FilteredData),
	}
}

func (app *App) processStartAndEndBlocks(line string, appconfig *AppConfig) {
	if !appconfig.hasStartBlockPattern && !appconfig.hasEndBlockPattern {
		app.clientValuesMap = make(map[string]*FilteredData)
		return
	}
	// TODO, implement start pattern and end pattern logic
	// TODO, reset valuesMap on start and end Blocks
}

func (app *App) filterData(line string, appconfig *AppConfig, logconfig *LogLineConfig) {

	elementsFound := false
	app.processStartAndEndBlocks(line, appconfig)

	for elementKey, ele := range logconfig.Elements {
		text := ""

		startPat := ele.StartPattern
		startIdx := 0
		if startPat == gStartOfLine {
			startPat = ""
			startIdx = 0
		} else {
			startIdx = strings.Index(line, startPat)
		}
		// continue if the StartPattern is not found
		if startIdx < 0 {
			continue
		}

		if len(ele.EndPattern) > 0 {
			endPat := ""
			endIdx := len(line)
			if ele.EndPattern != gEndOfLine {
				endPat = ele.EndPattern
				idx := strings.Index(line[startIdx+len(startPat):], endPat)
				// continue if EndPattern is not found
				if idx < 0 {
					continue
				}
				endIdx = startIdx + len(startPat) + idx
			}
			text = line[startIdx+len(startPat) : endIdx]
		} else if ele.PatternLength > 0 {
			if ele.PatternLength >= len(line) {
				// return if PatternLength is more than the length of the line
				// TODO, write testcases for this
				continue
			}
			endIdx := startIdx + len(startPat) + ele.PatternLength
			text = line[startIdx+len(startPat) : endIdx]
		}

		if logconfig.TrimSpaces {
			text = strings.TrimSpace(text)
		}
		if !ele.AllowEmpty && len(text) == 0 {
			text = "N/F"
		}
		elementsFound = true
		app.clientValuesMap[elementKey] = &FilteredData{Text: text}
	}
	if !elementsFound {
		return
	}

	if app.interactive {
		// TODO, change this fmt to logger; In console application use stdout (no logfile); In test create logger for debugging
		fmt.Println(logconfig.FormattedLine(line))
		for _, ele := range logconfig.Elements {
			fmt.Println(ele.Formatted())
		}
		fileutils.ReadStdin()
	}
	// TODO, Make seaprate interface for passing a static and dynamic configs to the clients
	appconfig.ClientConfig.Tag = logconfig.Tag

	if appconfig.hasEndBlockPattern {
		// TODO, check if the end block pattern matches then write output csv values
		// TODO, write a testcase
		if algoutils.StringContainsAll(line, appconfig.EndBlockPattern) {
			app.clientCallback(appconfig.ClientConfig, app.clientValuesMap)
		}
	} else {
		app.clientCallback(appconfig.ClientConfig, app.clientValuesMap)
	}
}

// Run starts the application
func (app *App) Run(clientCallback ClientCallbackType) {
	if clientCallback == nil {
		fmt.Println("ERROR, clientCallback is nil. Please provide a not-nil callback")
		return
	}
	app.clientCallback = clientCallback
	config := NewConfig(app.configFile, app.anchorFiles)
	if config == nil {
		return
	}

	lpr := logparser.NewLogParser()

	lpr.AddFileSources(app.inputFile)

	for _, appconfig := range config.Apps {
		for _, logconfig := range appconfig.LogLines {
			// we need to make a copy of logconfig here otherwise same logconfig is passed to logparser lambda
			logConfigCopy := logconfig
			appconfigCopy := appconfig
			lpr.AddConfig(logparser.Config{
				Patterns: logConfigCopy.Patterns,
				OnEachLineFunc: func(c *logparser.OnEachLineConfig) {
					app.filterData(c.Line, appconfigCopy, logConfigCopy)
				},
			})
		}
	}
	lpr.Run()
}

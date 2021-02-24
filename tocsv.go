package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"tocsv/filterlogs"

	"github.com/parmaanu/goutils/errorutils"
	"github.com/parmaanu/goutils/fileutils"

	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

type appDataType struct {
	Header []*filterlogs.MetaInfoType
	Data   [][]string
}

// Tocsv stores an instance of tocsv
type Tocsv struct {
	AppData       map[string]*appDataType // key is AppName
	Logfilter     *filterlogs.App
	PrintOnStdout bool
	Config        *TocsvConfig
}

const (
	gTagColumnName = "__tag__"
)

// NewTocsv returns a new Tocsv instance
func NewTocsv(inputFiles []string, configFile string, anchorFiles []string, printOnStdout, interactiveMode bool) *Tocsv {
	tocsvConfig := NewToCsvConfig(configFile)
	if tocsvConfig == nil {
		return nil
	}

	if len(anchorFiles) == 0 && len(tocsvConfig.AnchorFiles) > 0 {
		anchorFiles = tocsvConfig.AnchorFiles
	}

	// create LogDirectory if it does not exist
	if !printOnStdout {
		if len(tocsvConfig.LogDirectory) == 0 {
			tocsvConfig.LogDirectory = "."
		}
		absDir, err := tilde.Expand(tocsvConfig.LogDirectory)
		errorutils.PanicOnErr(err)
		if !fileutils.DirExist(absDir) {
			err = os.Mkdir(absDir, 0755)
			errorutils.PanicOnErr(err)
		}
		tocsvConfig.LogDirectory = absDir
	}

	return &Tocsv{
		AppData:       make(map[string]*appDataType),
		Logfilter:     filterlogs.NewApp(inputFiles, configFile, anchorFiles, interactiveMode),
		PrintOnStdout: printOnStdout,
		Config:        tocsvConfig,
	}
}

// Run runs the tocsv and prints output as csv
func (a *Tocsv) Run() {
	a.Logfilter.Run(a.callback)
	a.writeCsv()
}

func (a *Tocsv) callback(config *filterlogs.ClientConfigType, filteredDataMap map[string]*filterlogs.FilteredData) {
	appData, appExists := a.AppData[config.AppName]
	if !appExists {
		a.AppData[config.AppName] = &appDataType{}
		appData = a.AppData[config.AppName]
		appData.Header = config.MetaInfo
	}

	// TODO, avoid this hand-crafting of a csv file. Use something like dataframe in golang. This is error prone
	outputRecord := []string{}
	// Add tag in outputRecord if PrintTagInOutput is true in config
	if a.Config.PrintTagInOutput {
		outputRecord = append(outputRecord, config.Tag)
	}

	for _, metaInfo := range appData.Header {
		eleKey := metaInfo.ElementKey
		if filteredData, eleExists := filteredDataMap[eleKey]; eleExists {
			if strings.Index(filteredData.Text, ",") > -1 {
				outputRecord = append(outputRecord, "\""+filteredData.Text+"\"")
			} else {

				outputRecord = append(outputRecord, filteredData.Text)
			}
		} else {
			outputRecord = append(outputRecord, "N/A")
		}
	}
	appData.Data = append(appData.Data, outputRecord)
}

func (a *Tocsv) writeCsv() {

	for appName, appData := range a.AppData {
		// create header
		header := []string{}
		if a.Config.PrintTagInOutput {
			header = []string{gTagColumnName}
		}
		for _, metaInfo := range appData.Header {
			if len(metaInfo.ColumnName) > 0 {
				header = append(header, metaInfo.ColumnName)
			} else {
				header = append(header, metaInfo.ElementKey)
			}
		}

		if a.PrintOnStdout {
			printAsCsv(appName, os.Stdout, header, appData.Data)
		} else {
			dt := time.Now().Format("20060201") // YYYYMMDD format
			fname := fmt.Sprintf("%s/%s.%s.csv", a.Config.LogDirectory, appName, dt)
			oFile, err := os.Create(fname)
			if errorutils.PrintOnErr("ERROR while opening file for writing, "+fname, err) {
				continue
			}
			defer oFile.Close()

			printAsCsv(appName, oFile, header, appData.Data)
			fmt.Println("Data fetched in", fname)
		}
	}
}

func printAsCsv(appName string, w io.Writer, header []string, data [][]string) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	err := writer.Write(header)
	if errorutils.PrintOnErr("ERROR while writing header to output for app: "+appName, err) {
		return
	}
	for _, record := range data {
		err = writer.Write(record)
		if errorutils.PrintOnErr("ERROR while writing record to output for app: "+appName, err) {
			continue
		}
	}
}

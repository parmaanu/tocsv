package filterlogs

import (
	"fmt"
	"github.com/parmaanu/goutils/algoutils"
	"github.com/parmaanu/goutils/errorutils"
	"github.com/parmaanu/goutils/filesystem"
	"github.com/parmaanu/goutils/fileutils"
	"github.com/parmaanu/goutils/findutils"
	"os"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/lithammer/dedent"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

// MetaInfoType stores the meta information about the application like ColumnName, ElementKey
type MetaInfoType struct {
	ElementKey string
	ColumnName string
}

// ClientConfigType contains the config for client applications after extracting the pattern.
type ClientConfigType struct {
	AppName  string
	Tag      string // this tag changes for each matched logline // TODO, find a way to find out a better way to pass on dynamic config for each logline (or log block)
	MetaInfo []*MetaInfoType
}

// AppConfig stores the configuration for an individual app
type AppConfig struct {
	AppName           string           `yaml:"AppName"`
	StartBlockPattern []string         `yaml:"StartBlockPattern"`
	EndBlockPattern   []string         `yaml:"EndBlockPattern"`
	OutputElements    []string         `yaml:"OutputElements"`
	LogLines          []*LogLineConfig `yaml:"LogLines"`

	hasStartBlockPattern bool
	hasEndBlockPattern   bool
	ClientConfig         *ClientConfigType
}

// Config is the main application config
type Config struct {
	Apps []*AppConfig `yaml:"Apps"`
}

// NewConfig returns a config instance after reading configFiles
func NewConfig(configFile string, anchorFiles []string) *Config {
	absAnchorFiles := []string{}
	for _, fname := range anchorFiles {
		absfname, _ := tilde.Expand(fname)
		if !fileutils.FileExist(absfname) {
			fmt.Println("WARN:", absfname, "does not exist")
			continue
		}
		absAnchorFiles = append(absAnchorFiles, absfname)
	}
	absfname, _ := tilde.Expand(configFile)
	if !fileutils.FileExist(absfname) {
		fmt.Println("ERROR:", absfname, "does not exist")
		os.Exit(1)
	}
	reader, err := filesystem.Open(absfname)
	errorutils.PanicOnErr(err)
	// TODO, think about if I need to close it or not?
	// Anirudh: Close should be a member function of filesystem but the file interface
	// defer fs.Close()

	decoder := yaml.NewDecoder(reader, yaml.ReferenceFiles(absAnchorFiles...), yaml.UseOrderedMap())
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		errorutils.PrintOnErr("ERROR, while decoding config", err)
		return nil
	}

	if !config.Verify() {
		return nil
	}

	if !config.readAndStoreSortedElementKeys(absfname, absAnchorFiles) {
		return nil
	}
	return config
}

// Verify verifies the config file
func (config *Config) Verify() bool {
	return config.verifyAppConfig() && config.verifyLoglineConfig()
}

func (config *Config) verifyAppConfig() bool {
	if len(config.Apps) == 0 {
		fmt.Println("No apps configured in the config")
		return false
	}

	for _, appconfig := range config.Apps {
		if len(appconfig.AppName) == 0 {
			fmt.Println("AppName cannot be empty ", appconfig)
			return false
		}
		// appconfig.ClientConfig = &ClientConfigType{AppName: appconfig.AppName}
		if len(appconfig.LogLines) == 0 {
			fmt.Println("No LogLines found in the config. Please provide LogLines config")
			return false
		}

		for _, logline := range appconfig.LogLines {
			if len(logline.Tag) == 0 {
				fmt.Println("Please provide a Tag in logline config", logline)
				return false
			}
			if len(logline.Patterns) == 0 {
				fmt.Println("Please provide Patterns to filter logs in logline config", logline)
				return false
			}
			if len(logline.ExampleLine) == 0 {
				fmt.Println("Please provide ExampleLine to remember the corresponding logline", logline)
				return false
			}
			if !algoutils.StringContainsAll(logline.ExampleLine, logline.Patterns) {
				fmt.Println("Config patterns not found in the example line, patterns: ", logline.Patterns, logline.Tag, logline.ExampleLine)
				return false
			}

			if len(logline.Elements) == 0 {
				exampleElements := dedent.Dedent(`
			Elements:
			  TimestampKey:
			    ColumnName: timestamp
				StartPattern: '^'
				PatternLength: 26
			  InstrumentIdKey:
			    ColumnName: instrumentId
				StartPattern: 'instrumentId: '
				EndPattern: ','
				AllowEmpty: true
			`)
				fmt.Println("Please provide Elements to be printed in output csv file", logline, ". Following is an example:", exampleElements)
				return false
			}
		} // LogLines loop
	}
	return true
}

func (config *Config) verifyLoglineConfig() bool {
	for _, appconfig := range config.Apps {
		elementKeys := []string{}

		for _, logline := range appconfig.LogLines {
			columnNames := []string{}

			for eleKey, ele := range logline.Elements {

				appconfig.hasStartBlockPattern = len(appconfig.StartBlockPattern) > 0 && len(appconfig.StartBlockPattern[0]) > 0
				appconfig.hasEndBlockPattern = len(appconfig.EndBlockPattern) > 0 && len(appconfig.EndBlockPattern[0]) > 0
				if appconfig.hasStartBlockPattern && appconfig.hasEndBlockPattern {
					// check for repeatative element key across different log lines
					if findutils.ContainsString(elementKeys, eleKey) {
						fmt.Println("Repeated ElementKeys found in LogLines config, Please provide unique element keys within an app", logline.Tag, eleKey)
						return false
					}

					// check unique column name in each log line config, ColumnName is not compulsory
					if len(ele.ColumnName) > 0 && findutils.ContainsString(columnNames, ele.ColumnName) {
						fmt.Println("Repeated ColumnName in Elements config, Please provide unique set of column names for each logLine", logline.Tag, eleKey, ele)
						return false
					}
					elementKeys = append(elementKeys, eleKey)
					columnNames = append(columnNames, ele.ColumnName)
				}

				if ele.StartPattern != gStartOfLine && !strings.Contains(logline.ExampleLine, ele.StartPattern) {
					fmt.Println("StartPattern not found in ExampleLine, StartPattern:", ele.StartPattern, logline.Tag, eleKey, ele, logline.ExampleLine)
					return false
				}
				if len(ele.EndPattern) > 0 && ele.EndPattern != gEndOfLine && !strings.Contains(logline.ExampleLine, ele.EndPattern) {
					fmt.Println("EndPattern not found in ExampleLine, EndPattern:", ele.EndPattern, logline.Tag, eleKey, ele, logline.ExampleLine)
					return false
				}
				if len(ele.StartPattern) == 0 {
					fmt.Println("Please provide the StartPattern in the ElementConfig", logline.Tag, eleKey, ele)
					return false
				}
				if ele.PatternLength < 0 {
					fmt.Println("Negative PatternLength is not supported", logline.Tag, eleKey, ele)
					return false
				}
				if ele.PatternLength > 0 && len(ele.EndPattern) > 0 {
					fmt.Println("Either provide PatternLength or EndPattern, simulaneously both are not supported", logline.Tag, eleKey, ele)
					return false
				}
				if ele.StartPattern == gEndOfLine || ele.EndPattern == gStartOfLine {
					fmt.Println("Please provide correct startPattern and endPattern", logline.Tag, eleKey, ele)
					return false
				}
			}
		}
	}
	return true
}

func (config *Config) readAndStoreSortedElementKeys(absConfigFile string, absAnchorFiles []string) bool {
	reader, err := filesystem.Open(absConfigFile)
	errorutils.PanicOnErr(err)

	type tempLogLineConfig struct {
		Tag              string        `yaml:"Tag"`
		ElementsMapSlice yaml.MapSlice `yaml:"Elements"`
	}

	type tempAppConfig struct {
		AppName  string               `yaml:"AppName"`
		LogLines []*tempLogLineConfig `yaml:"LogLines"`
	}

	type tempConfig struct {
		App []*tempAppConfig `yaml:"Apps"`
	}
	tempDecoder := yaml.NewDecoder(reader, yaml.ReferenceFiles(absAnchorFiles...), yaml.UseOrderedMap())
	tempC := &tempConfig{}
	if err := tempDecoder.Decode(tempC); err != nil {
		errorutils.PrintOnErr("ERROR", err)
		return false
	}

	for i, tempApp := range tempC.App {
		app := config.Apps[i]
		if app.AppName != tempApp.AppName {
			fmt.Println("ERROR:", "AppName does not match after reading again the temp config", app.AppName, tempApp.AppName)
			return false
		}
		app.ClientConfig = &ClientConfigType{AppName: app.AppName}
		elementKeys := map[string]bool{}
		for j, tempLogLineConfig := range tempApp.LogLines {
			loglineConfig := app.LogLines[j]
			if loglineConfig.Tag != tempLogLineConfig.Tag {
				fmt.Println("ERROR:", "Tag does not match after reading again the config", loglineConfig.Tag, tempLogLineConfig.Tag)
				return false
			}
			elements := loglineConfig.Elements
			for _, e := range tempLogLineConfig.ElementsMapSlice {
				eleKey, ok := e.Key.(string)
				if !ok {
					fmt.Println("ERROR:", "Cannot convert elementKey to string after reading again the config", e, tempLogLineConfig.Tag)
					return false
				}
				element, exists := elements[eleKey]
				if !exists {
					fmt.Println("ERROR:", "Cannot find elementKey in actual config Element map after reading again the config", eleKey, elements)
					return false
				}
				if _, alreadyExists := elementKeys[eleKey]; alreadyExists {
					continue
				}
				app.ClientConfig.MetaInfo = append(app.ClientConfig.MetaInfo, &MetaInfoType{ElementKey: eleKey, ColumnName: element.ColumnName})
				elementKeys[eleKey] = true
			}
		}
	}
	return true
}

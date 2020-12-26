// Package logparser exposes API related to reading of log files
// Sets up the multisource linereader and invoked callback on each matching line
package logparser

import (
	"fmt"

	"github.com/parmaanu/goutils/algoutils"
)

// LineParsingConfig is passed in LineParsingFunc which is called for every line
type LineParsingConfig struct {
	Line string
	Pat  string
}

// OnEachLineConfig is passed in OnEachLineFunc which is called after each line is parsed
// TODO: a better name is needed this is not a Config, but a result
type OnEachLineConfig struct {
	Line string
	Pat  string
	Data []string
}

// Config depicts the each element config and a specific action
type Config struct {
	Patterns        []string
	LineParsingFunc func(config *LineParsingConfig) []string
	OnEachLineFunc  func(config *OnEachLineConfig)
}

// LogParser is the main struct will contains the patterns and the different sources
type LogParser struct {
	mslr     MultiSourceLineReader
	patterns []*Config
}

func (lp *LogParser) processLine(line string) {
	for _, config := range lp.patterns {
		if !algoutils.StringContainsAll(line, config.Patterns) {
			continue
		}

		data := []string{}
		if config.LineParsingFunc != nil {
			// TODO, for backward compatibility we are passing on the first pattern, it should pass a string of patterns
			// and provide a comparing utility inside logparser
			data = config.LineParsingFunc(&LineParsingConfig{
				Line: line,
				Pat:  config.Patterns[0],
			})
		} else {
			data = []string{line}
		}

		if config.OnEachLineFunc != nil && len(data) > 0 && len(data[0]) > 0 {
			config.OnEachLineFunc(&OnEachLineConfig{
				Line: line,
				Pat:  config.Patterns[0],
				Data: data,
			})
		}
		// break if you found any pattern
		break
	}
}

// Run starts the processing of different sources
func (lp *LogParser) Run() {
	nextLine, err := lp.mslr.NextLine()
	for err != -1 {
		lp.processLine(nextLine)
		nextLine, err = lp.mslr.NextLine()
	}
}

// AddFileSources takes a list of filenames and created FileLineReader for
// reading from those files
func (lp *LogParser) AddFileSources(filenames ...string) {
	for _, filename := range filenames {
		flr, err := NewFileLineReader(filename)
		if err == nil {
			lp.mslr.AddSources(flr)
		} else {
			fmt.Println("Failed to create FLR err ", err)
		}
	}
}

// AddConfig register the current config to list of patterns the
// LogParser is interested in
func (lp *LogParser) AddConfig(config Config) bool {
	if len(config.Patterns) == 0 || len(config.Patterns[0]) == 0 {
		return false
	}
	lp.patterns = append(lp.patterns, &config)
	return true
}

// NewLogParser creats an instance of LogParser
func NewLogParser() *LogParser {
	return &LogParser{}
}

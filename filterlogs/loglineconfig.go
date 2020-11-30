package filterlogs

import (
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/muesli/reflow/wordwrap"
	"golang.org/x/crypto/ssh/terminal"
)

// LogLineConfig stores the configuration for each logline
type LogLineConfig struct {
	Tag         string                    `yaml:"Tag"`
	TrimSpaces  bool                      `yaml:"TrimSpaces"`
	Patterns    []string                  `yaml:"Patterns"`
	ExampleLine string                    `yaml:"ExampleLine"`
	Elements    map[string]*ElementConfig `yaml:"Elements"`

	// TODO, later on we can rename Elements with Columns and ColumnConfig if required. It is also possible that each
	// element does not result in a column
	cachedFormattedLineConfig string
}

// FormattedExampleLine returns the example log line formatted with config patterns
func (logline *LogLineConfig) FormattedExampleLine() string {
	if len(logline.cachedFormattedLineConfig) > 0 {
		return logline.cachedFormattedLineConfig
	}
	logline.cachedFormattedLineConfig = logline.FormattedLine(logline.ExampleLine)
	return logline.cachedFormattedLineConfig
}

// FormattedLine returns the example log line formatted with config patterns. This can be used for debugging purposes to
// figure out what text has been filtered out by the provided logline config.
func (logline *LogLineConfig) FormattedLine(inputLine string) string {
	red := promptui.Styler(promptui.FGRed)
	green := promptui.Styler(promptui.FGGreen, promptui.FGBold)

	line := inputLine
	for _, ele := range logline.Elements {

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
				if endIdx < 0 {
					continue
				}
				endIdx = startIdx + idx + len(startPat)
			}
			result := line[startIdx+len(startPat) : endIdx]
			line = line[:startIdx] + red(startPat) + green(result) + red(endPat) + line[endIdx+len(endPat):]
		} else if ele.PatternLength > 0 {
			if ele.PatternLength >= len(line) {
				// return if PatternLength is more than the length of the line
				// TODO, write testcases for this
				continue
			}
			endIdx := startIdx + len(startPat) + ele.PatternLength + 1
			result := line[startIdx+len(startPat) : endIdx]
			line = line[:startIdx] + red(startPat) + green(result) + line[endIdx:]
		}
	}

	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		line = wordwrap.String(line, width-1)
	}
	return line
}

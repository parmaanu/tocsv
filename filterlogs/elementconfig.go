package filterlogs

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// ElementConfig store the config for each pattern in a line
type ElementConfig struct {
	ColumnName    string `yaml:"ColumnName"`
	StartPattern  string `yaml:"StartPattern"`
	EndPattern    string `yaml:"EndPattern"`
	AllowEmpty    bool   `yaml:"AllowEmpty"` // When this is set as true then empty values does not print N/F for this column
	PatternLength int    `yaml:"PatternLength"`
	After         string `yaml:"After"`

	cacheFormattedConfig string
}

// Formatted returns the representation of ElementConfig in one single line string
func (ele *ElementConfig) Formatted() string {
	if len(ele.cacheFormattedConfig) > 0 {
		return ele.cacheFormattedConfig
	}
	output := []string{}
	blue := promptui.Styler(promptui.FGBlue, promptui.FGBold)

	if len(ele.ColumnName) > 0 {
		output = append(output, fmt.Sprintf("%s: %-23s", blue("ColName"), ele.ColumnName))
	}
	output = append(output, fmt.Sprintf("%s: %-23s", blue("Start"), "'"+ele.StartPattern+"'"))
	if len(ele.EndPattern) > 0 {
		output = append(output, fmt.Sprintf("%s: %s", blue("End"), "'"+ele.EndPattern+"'"))
	}
	if ele.PatternLength > 0 {
		output = append(output, fmt.Sprintf("%s: %d", blue("Len"), ele.PatternLength))
	}
	ele.cacheFormattedConfig = strings.Join(output, " ")
	return ele.cacheFormattedConfig
}

package filterlogs

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	quit = "Quit"
)

// TODO, remove if not required
// // DumpSampleConfig dumps a sample config on terminal
// func DumpSampleConfig() {
//
// }

// PrintInteractiveConfig reads and print filterlogs config interactively
func PrintInteractiveConfig(configFile string, anchorFiles []string) {
	config := NewConfig(configFile, anchorFiles)
	if config == nil {
		return
	}
	showConfig(config)
}

func showConfig(config *Config) {
	apps := append(config.Apps, &AppConfig{AppName: quit})

	searcher := func(input string, index int) bool {

		app := apps[index]
		name := strings.Replace(strings.ToLower(app.AppName), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	templates := &promptui.SelectTemplates{
		// Label:    "{{ . }}?",
		Label:    "",
		Active:   "> {{ .AppName | cyan | bold }}",
		Inactive: "  {{ .AppName | cyan }}",
		Selected: "> {{ .AppName | green }}",
		Details: `{{if eq .AppName "Quit"}} Quits the prompt!!!{{ else }}
--------- Configurations ----------
{{ printf "%-10v" "AppName:" | faint }}   {{ .AppName | bold }}
{{ printf "%-25v" "LogLines:" | faint -}}
{{ range .LogLines }}
{{ .FormattedExampleLine }}
  {{ printf "%-25v" "Tag" | faint }}: {{ .Tag }}
  {{ printf "%-25v" "Patterns" | faint }}: {{ range .Patterns }}"{{.}}" {{end}}
  {{ range $k, $v := .Elements }} 
	  {{- printf "- %-23v" $k | faint}}: {{ $v.Formatted }}
  {{end}}
{{- end}}
{{- end}}`,
	}

	prompt := promptui.Select{
		Label:     "Select an app to show its config",
		Items:     apps,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}
	_, selectedApp, err := prompt.Run()
	if err != nil {
		fmt.Println("Error while selecting an app", err.Error(), ", selected: ", selectedApp)
		return
	}
}

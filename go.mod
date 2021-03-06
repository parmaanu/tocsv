module tocsv

go 1.12

require (
	github.com/AlecAivazis/survey/v2 v2.2.8
	github.com/goccy/go-yaml v1.8.3
	github.com/google/go-cmp v0.3.0
	github.com/lithammer/dedent v1.1.0
	github.com/manifoldco/promptui v0.8.0
	github.com/muesli/reflow v0.2.0
	github.com/parmaanu/goutils v0.0.0-20201130155100-92dcaa7f6188
	github.com/parmaanu/showcsv v0.0.0-20201226140506-2d72b643f8de
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5
	golang.org/x/sys v0.0.0-20201017003518-b09fb700fbb7 // indirect
	gopkg.in/mattes/go-expand-tilde.v1 v1.0.0-20150330173918-cb884138e64c
)

replace github.com/parmaanu/goutils => ../goutils

replace github.com/parmaanu/showcsv => ../showcsv

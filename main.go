package main

import (
	"errors"
	"fmt"
	"tocsv/filterlogs"
	"tocsv/tocsvgo"

	"github.com/spf13/cobra"
)

var appname = "tocsv"

func main() {
	inputFiles := []string{}
	configFile := ""
	anchorFiles := []string{}
	printOnStdout := false
	printLogLines := false
	interactiveMode := false
	dumpConfig := false

	rootCmd := &cobra.Command{
		Use: appname,
		Run: func(cmd *cobra.Command, args []string) {
			if len(inputFiles) == 0 {
				// Read config and show config on stdout interactively
				tocsvConfig := tocsvgo.NewToCsvConfig(configFile)
				if tocsvConfig != nil {
					if len(anchorFiles) == 0 && len(tocsvConfig.AnchorFiles) > 0 {
						anchorFiles = tocsvConfig.AnchorFiles
					}
					filterlogs.PrintInteractiveConfig(configFile, anchorFiles)
				}
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(inputFiles) != 0 {
				return nil
			}
			if len(args) == 0 {
				return errors.New("please provide log file")
			}
			inputFiles = args[:]
			return nil
		},
	}

	defaultConfigFile := "~/." + appname + ".yaml"
	// TODO, current only one main config file is supported, rest are anchor files. Support multiple config files later
	// rootCmd.Flags().StringArrayVarP(&configFiles, "config", "c", []string{defaultConfigFile}, "input config yamls for tocsv app")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", defaultConfigFile, "input config yamls for tocsv app")
	rootCmd.Flags().StringArrayVarP(&anchorFiles, "anchor", "a", []string{}, "input anchor config yamls files")
	rootCmd.Flags().StringArrayVarP(&inputFiles, "files", "f", []string{}, "input logfiles for tocsv app")
	rootCmd.Flags().BoolVarP(&printOnStdout, "print", "p", false, "print the output on stdout")
	// TODO, printLogLines is a part of tocsv not filterlogs
	rootCmd.Flags().BoolVarP(&printLogLines, "loglines", "l", false, "log actual loglines with csv")
	rootCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "interactive mode on")
	// TODO, dump config
	rootCmd.Flags().BoolVarP(&dumpConfig, "dump-config", "d", false, "dump sample config")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}

	if len(inputFiles) == 0 {
		return
	}

	if len(inputFiles) > 1 {
		fmt.Println("ERROR, more than one logfile is not supported currently: ", inputFiles)
		return
	}

	tocsv := tocsvgo.NewTocsv(inputFiles, configFile, anchorFiles, printOnStdout, interactiveMode)
	if tocsv != nil {
		tocsv.Run()
	}
	tocsv.DisplayFetchedCsvs()
}

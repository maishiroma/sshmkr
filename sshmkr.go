package main

import (
	"fmt"
	"os"
	"os/user"
	"flag"	
	"sshmkr/help"
	"sshmkr/reader"
	"sshmkr/input"
	"sshmkr/commands"
)

//// Global Variables

// Flag Values
var helpFlagValue bool
var versionFlagValue bool
var configFlagValue string

// Initializes Program
func init() {
	flag.BoolVar(&helpFlagValue, "help", false, "help flag")
	flag.BoolVar(&helpFlagValue, "h", false, "help flag")

	flag.BoolVar(&versionFlagValue, "version", false, "version flag")
	flag.BoolVar(&versionFlagValue, "v", false, "version flag")

	usr , err := user.Current()
	if err != nil {
		fmt.Println("Error! Program cannot get current user!")
		os.Exit(1)
    }
	defaultConfigPath := fmt.Sprint(usr.HomeDir, "/.ssh/config")

	flag.StringVar(&configFlagValue, "path", defaultConfigPath, "Directory of ssh config")
	flag.StringVar(&configFlagValue, "p", defaultConfigPath, "Directory of ssh config")
}

// Main Execution of Program
func main() {

	configFile, configFileContents, configFileDecoded := sshmkr_reader.ParseConfigFile(configFlagValue)
	configTemplateFile, _, configTemplateFileDecoded := sshmkr_reader.ParseConfigFile(fmt.Sprintf("%s_templates", configFlagValue))
	defer configFile.Close()
	defer configTemplateFile.Close()

	// Setting up the subcommands and their flags
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addSource := addCmd.String("source", "", "Name of source template config to leverage")
	sshmkr_help.SetHelpContext(addCmd, "add")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteSource := deleteCmd.String("source", "", "Name of host config to remove")
	sshmkr_help.SetHelpContext(deleteCmd, "delete")

	copyCmd := flag.NewFlagSet("copy", flag.ExitOnError)
	copySource := copyCmd.String("source", "", "Name of host config to use as basis")
	sshmkr_help.SetHelpContext(copyCmd, "copy")

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	showSource := showCmd.String("source", "", "Name of host config to show")
	sshmkr_help.SetHelpContext(showCmd, "show")

	commentCmd := flag.NewFlagSet("comment", flag.ExitOnError)
	commentSource := commentCmd.String("source", "", "Name of host config to interact")
	sshmkr_help.SetHelpContext(commentCmd, "comment")

	if len(os.Args) < 2 {
		fmt.Println("Error! Expecting another argument: [add, delete, copy, show, list]")
		os.Exit(1)
	}

	flag.Parse()
	switch os.Args[1] {
		case "add":
			addCmd.Parse(os.Args[2:])

			template := sshmkr_reader.ReadSpecificTemplate(*addSource, configTemplateFileDecoded)
			headers := sshmkr_reader.ParseConfigHeaders(configFileContents)
			mainHeader, subHeader := sshmkr_input.SelectNewConfigLoc(headers)
			userAddedConfig := sshmkr_input.InterpolateUserInput(template)
			newOutput := sshmkr_commands.AddTemplatedConfig(mainHeader, subHeader, userAddedConfig, configFileContents)
			sshmkr_reader.WriteToConfigFile(configFlagValue, newOutput)

			fmt.Println("Sucessfully added new config to ssh_config!")
		case "delete":
			deleteCmd.Parse(os.Args[2:])

			newOutput := sshmkr_commands.RemoveHostConfig(*deleteSource, configFileContents)
			sshmkr_reader.WriteToConfigFile(configFlagValue, newOutput)
			
			fmt.Println("Sucessfully removed config from ssh_config!")
		case "copy":
			copyCmd.Parse(os.Args[2:])

			template := sshmkr_reader.ReadSpecificTemplate(*copySource, configFileDecoded)
			headers := sshmkr_reader.ParseConfigHeaders(configFileContents)
			mainHeader, subHeader := sshmkr_input.SelectNewConfigLoc(headers)
			userAddedConfig := sshmkr_input.InterpolateUserInput(template)
			newOutput := sshmkr_commands.AddTemplatedConfig(mainHeader, subHeader, userAddedConfig, configFileContents)
			sshmkr_reader.WriteToConfigFile(configFlagValue, newOutput)			

			fmt.Println("Sucessfuly copied one config to another config!")
		case "show":
			showCmd.Parse(os.Args[2:])

			sshmkr_commands.GetSpecificHostConfig(*showSource, configFileDecoded)
		case "comment":
			commentCmd.Parse(os.Args[2:])

			newOutput, hasCommented := sshmkr_commands.CommentHostConfig(*commentSource, configFileContents)
			sshmkr_reader.WriteToConfigFile(configFlagValue, newOutput)
			
			if hasCommented {
				fmt.Println("Sucessfully commented out config!")
			} else {
				fmt.Println("Sucessfully uncommented out config!")
			}
		default:
			if helpFlagValue == true {
				sshmkr_help.DefaultHelp()
			} else if versionFlagValue == true {
				sshmkr_help.PrintVersion()
			} else {
				fmt.Printf("Subcommand '%s' invalid. Available commands are: [add, delete, copy, show, list]\n", os.Args[1])
				os.Exit(1)
			}
	}
}
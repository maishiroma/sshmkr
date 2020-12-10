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

//// Program Functions

// Adds a new host config to a config file
// Returns the new config file contents
func addTemplatedConfig(mainHeader string, subHeader string, templateString string, fileContents []byte) string {
	/*
	*	The logic behind this is that we are adding in new config based on a passed template.
	* 	The user will pass in three flags  (two being config headers) and the name of the template used.
	*	The user will be prompted to enter in new values into the hostname (or can hit enter to use the )
	*	default). Once done, the new config will be put into the config file.
	*
	*/
	
	foundHeader := false
	foundSubHeader := false
	fileContentsArray := strings.Split(string(fileContents), "\n")
	
	for currIndex, currLine := range fileContentsArray {
		if (currLine == mainHeader) {
			foundHeader = true
		} else if (foundHeader == true && currLine == subHeader) {
			foundSubHeader = true
		}

		if foundHeader == true && foundSubHeader == true {
			if strings.Contains(currLine, sshmkr_reader.SUB_HEADER_IND) {
				// This means the previous line is the line that we are interested in replacing
				fileContentsArray[currIndex-1] = templateString
				break
			} else if currIndex + 1 >= len(fileContentsArray) {
				// We reached the end of the file, so we just place the new contents here
				fileContentsArray[currIndex] = templateString
				break
			}
		} 
	}

	newContents := strings.Join(fileContentsArray, "\n")
	return newContents
}

// Removes a specified host config from the ssh_config value
// And returns the updated file content
func removeHostConfig(hostname string, fileContents []byte) string {
	if len(hostname) <= 0 {
		fmt.Println("Source flag is empty! Please pass in a valid hostname to remove!")
		os.Exit(-1)
	}
	
	foundHost := false		// Has the method found the matching hostname
	finishDelete := false	// Has the method finished parsing out the specified hostname's config
	newFileIndex := 0		// The current array index that the new byte array is currently on
	numbLinesRemoved := 0	// How many lines will be omitted from the new file array

	hostSearch := fmt.Sprintf("Host %s", hostname)

	// In order to effectively delete a specific spot in the file
	// We take the original file contents and make a copy of it to another array
	// without copying over the specified config
	fileContentLines := strings.Split(string(fileContents), "\n")
	newFileContents := make([]string, len(fileContentLines))
	
	for currIndex, currLine := range fileContentLines {
		
		if strings.Contains(currLine, hostSearch) {
			// We found the specified hostname in the original file
			foundHost = true
		} 

		if foundHost == true && finishDelete == false {
			// Once we find the matching host, we omit it from the copy
			numbLinesRemoved = numbLinesRemoved + 1
			if currLine == "" || currIndex + 1 >= len(fileContentLines) {
				// We reached the end of the config we wanted to remove and resume copying
				finishDelete = true
			}
		} else {
			// We copy the original file to the new file array
			newFileContents[newFileIndex] = currLine
			newFileIndex = newFileIndex + 1
		}
	}

	if foundHost == false {
		fmt.Println("Cannot find specified hostname in config. Typo maybe?")
		os.Exit(-1)
	}

	// We take a slice of the new array (since we are removing X lines from the file)
	// and return it
	newFileArray := strings.Join(newFileContents, "\n")
	return newFileArray[:len(newFileArray) - numbLinesRemoved]
}

// Comments/Uncomments a specific host config depending if it was already commented or not
// Return the updated file contents and if it did comment it out
func commentHostConfig(hostname string, fileContents []byte) (string, bool) {
	if len(hostname) <= 0 {
		fmt.Println("Source flag is empty! Please pass in a valid hostname to remove!")
		os.Exit(-1)
	}
	
	foundHost := false		// Has the method found the matching hostname
	hasCommented := false	// Did we comment out the config?

	hostSearch := fmt.Sprintf("Host %s", hostname)
	fileContentLines := strings.Split(string(fileContents), "\n")
	
	for currIndex, currLine := range fileContentLines {
		
		if foundHost == false {
			if strings.Contains(currLine, hostSearch) {
				// We found the specified hostname in the original file
				foundHost = true
			} 
		} 
		
		if foundHost == true {
			if currLine == "" || currIndex + 1 >= len(fileContentLines) {
				// We reached the end of the modding and break out of the loop
				break
			} else if !strings.Contains(currLine, sshmkr_reader.COMMENT_IND) {
				// Once we find the matching host, we omit/readd it in the config
				fileContentLines[currIndex] = fmt.Sprintf("%s%s", sshmkr_reader.COMMENT_IND, currLine)
				hasCommented = true
			} else {
				fileContentLines[currIndex] = currLine[len(sshmkr_reader.COMMENT_IND):]
			}	
		}		
	}

	if foundHost == false {
		fmt.Println("Cannot find specified hostname in config. Typo maybe?")
		os.Exit(-1)
	}

	return strings.Join(fileContentLines, "\n"), hasCommented
}

// Prints out a specific host configuration out to standard output
func getSpecificHostConfig(hostname string, parsedConfig *ssh_config.Config) {
	for _, host := range parsedConfig.Hosts {
		if host.Patterns[0].String() == hostname {
			// Once we found the desired host, we print it out in its entirety and leave the statement
			fmt.Println("Host ", hostname)
			for _, node := range host.Nodes {
				output := node.String()

				if sshmkr_reader.CheckIfValid(output) {
					fmt.Println(output)
				}
			}
			break
		}
	}
}

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
			mainHeader, subHeader := selectNewConfigLoc(headers)
			userAddedConfig := interpolateUserInput(template)
			newOutput := addTemplatedConfig(mainHeader, subHeader, userAddedConfig, configFileContents)
			sshmkr_reader.WriteToConfigFile(configFlagValue, newOutput)

			fmt.Println("Sucessfully added new config to ssh_config!")
		case "delete":
			deleteCmd.Parse(os.Args[2:])

			newOutput := removeHostConfig(*deleteSource, configFileContents)
			sshmkr_reader.WriteToConfigFile(configFlagValue, newOutput)
			
			fmt.Println("Sucessfully removed config from ssh_config!")
		case "copy":
			copyCmd.Parse(os.Args[2:])

			template := sshmkr_reader.ReadSpecificTemplate(*copySource, configFileDecoded)
			headers := sshmkr_reader.ParseConfigHeaders(configFileContents)
			mainHeader, subHeader := selectNewConfigLoc(headers)
			userAddedConfig := interpolateUserInput(template)
			newOutput := addTemplatedConfig(mainHeader, subHeader, userAddedConfig, configFileContents)
			sshmkr_reader.WriteToConfigFile(configFlagValue, newOutput)			

			fmt.Println("Sucessfuly copied one config to another config!")
		case "show":
			showCmd.Parse(os.Args[2:])

			getSpecificHostConfig(*showSource, configFileDecoded)
		case "comment":
			commentCmd.Parse(os.Args[2:])

			newOutput, hasCommented := commentHostConfig(*commentSource, configFileContents)
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
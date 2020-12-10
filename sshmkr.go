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

// Parses the passed config file location to the program
// Returns the open file, the file contents, and the decoded config file
func parseConfigFile(configLoc string) (*os.File, []byte ,*ssh_config.Config) {
	config_file, err1 := os.Open(configLoc)
	config_file_contents, err2 := ioutil.ReadFile(configLoc)
	config_file_decoded, err3 := ssh_config.Decode(config_file)

	if err1 != nil {
		fmt.Println("Error! The config location: ", configLoc, " cannot be opned!")
		os.Exit(1)
	} else if err2 != nil {
		fmt.Println("Error! The config location: ", configLoc, " cannot be read!")
		os.Exit(1)
	}else if err3 != nil {
		fmt.Println("Error! The config location: ", configLoc, " cannot be parsed!")
		os.Exit(1)
	}

	return config_file, config_file_contents, config_file_decoded
}

// Writes out the passed in string into a new file
func writeToConfigFile(configLoc string, fileContents string) {
	err := ioutil.WriteFile(configLoc, []byte(fileContents), 644)
	if err != nil {
		fmt.Println("Error! The location", configLoc, " cannot be written!")
		os.Exit(1)
	}
}

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
			if strings.Contains(currLine, SUB_HEADER_IND) {
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
			} else if !strings.Contains(currLine, COMMENT_IND) {
				// Once we find the matching host, we omit/readd it in the config
				fileContentLines[currIndex] = fmt.Sprintf("%s%s", COMMENT_IND, currLine)
				hasCommented = true
			} else {
				fileContentLines[currIndex] = currLine[len(COMMENT_IND):]
			}	
		}		
	}

	if foundHost == false {
		fmt.Println("Cannot find specified hostname in config. Typo maybe?")
		os.Exit(-1)
	}

	return strings.Join(fileContentLines, "\n"), hasCommented
}

// Returns a ConfigTemplate object that contains information on a given template
func readSpecificTemplate(hostname string, config_template *ssh_config.Config) sshmkr_templates.ConfigTemplate {
	for _, host := range config_template.Hosts {
		if  checkIfExistingHostname(hostname, host.Patterns[0].String()) {
			currIndex := 0
			formatted_template_string := "\nHost %s\n"
			template_kv := make([]ssh_config.KV, len(host.Nodes))
			
			template_kv[currIndex] = ssh_config.KV{Key: "Host", Value: "NewHost", Comment: ""}
			currIndex = currIndex + 1
			for _, node := range host.Nodes {
				nodeRendered := node.String()
				
				if checkIfValid(nodeRendered) {
					// Key Index = 0; Value Index = 1
					kvPair := strings.Fields(nodeRendered)
					formatted_template_string = formatted_template_string + "\t" + kvPair[0] + " %s \n"

					// The order of the interpolation in the format template string also correlates
					// to the order of the values in the array
					template_kv[currIndex] = ssh_config.KV{Key: kvPair[0], Value: kvPair[1], Comment: ""}
					currIndex = currIndex + 1
				}
			}

			// We then create a struct object from the data we gathered and return it out
			// Note that the length of the default values is always the same number as the
			// numbrrt of special replacement chars
			return sshmkr_templates.ConfigTemplate{KeyPairs: template_kv[:currIndex], FormattedString: formatted_template_string}
		}
	}

	// Only comes here if the passed in template name does not match any existing ones
	fmt.Println("Cannot find template specified! Typo maybe?")
	os.Exit(-1)
	return sshmkr_templates.ConfigTemplate{}
}

// Prints out a specific host configuration out to standard output
func getSpecificHostConfig(hostname string, parsedConfig *ssh_config.Config) {
	for _, host := range parsedConfig.Hosts {
		if host.Patterns[0].String() == hostname {
			// Once we found the desired host, we print it out in its entirety and leave the statement
			fmt.Println("Host ", hostname)
			for _, node := range host.Nodes {
				output := node.String()

				if checkIfValid(output) {
					fmt.Println(output)
				}
			}
			break
		}
	}
}

// Takes in a templated string and user input to return a filled host config
func interpolateUserInput(template sshmkr_templates.ConfigTemplate) string {
	templateString := template.GetTemplatedString()
	for currIndex := 0; currIndex < template.GetNumKeyPairs(); currIndex = currIndex + 1 {
		kvPair := template.GetKeyPair(currIndex)
		var userInput string

		fmt.Printf("Enter a value for %s [ default: %s ]: ", kvPair[0], kvPair[1])
		fmt.Scanln(&userInput)
		if userInput == "" {
			templateString = strings.Replace(templateString, "%s", kvPair[1], 1)
		} else {
			templateString = strings.Replace(templateString, "%s", userInput, 1)
		}
		
	}
	return templateString
}

// Takes in a parsed ssh config and outputs all of the relevant header comments
// Returns an array of headerBlocks, which are logical groupings of ssh configs
func parseConfigHeaders(fileContents []byte) []sshmkr_templates.HeaderBlock {
	fileContentsArray := strings.Split(string(fileContents), "\n")
	
	headerBlocks := make([]sshmkr_templates.HeaderBlock, len(fileContentsArray))
	headerBlocksSize := 0

	newHeaderBlock := ""
	newSubHeadersBlock := ""	// sub headers will be stored in one string, formatted by \n
	lastHeader := ""			// Saves a ref to the last parsed header value
	
	for currIndex, currLine := range fileContentsArray {
		
		if currIndex + 1 == len(fileContentsArray) {
			// Once we reached the end of the file, we parse out the last remaining HeaderBlock
			newSubHeadersBlock = strings.TrimRight(newSubHeadersBlock, "\n")
			splitHeaders := strings.Split(newSubHeadersBlock, "\n")
			headerBlocks[headerBlocksSize] = sshmkr_templates.HeaderBlock{MainHeader: newHeaderBlock, SubHeaders: splitHeaders}
			headerBlocksSize = headerBlocksSize + 1
		} else {
			commentStartIndex := strings.Index(currLine, " ")
			if commentStartIndex != -1 {
				// We come in here once we are on a line that is formatted to be a header
				parsedLine := currLine[:commentStartIndex]

				if parsedLine == MAIN_HEADER_IND {
					if lastHeader == SUB_HEADER_IND {
						// If we reach a new header block, we save the HeaderBlock that we made
						// to the array and start making a new one
						newSubHeadersBlock = strings.TrimRight(newSubHeadersBlock, "\n")
						splitHeaders := strings.Split(newSubHeadersBlock, "\n")
						headerBlocks[headerBlocksSize] = sshmkr_templates.HeaderBlock{MainHeader: newHeaderBlock, SubHeaders: splitHeaders}
						headerBlocksSize = headerBlocksSize + 1
					}
					
					// We encounted a main header, making it the start of a new HeaderBlock obj
					newHeaderBlock = currLine
					newSubHeadersBlock = ""
					lastHeader = parsedLine
					
				} else if parsedLine == SUB_HEADER_IND && newHeaderBlock != "" {
					// We are adding more to the current HeaderBlock that we are making
					// As long as we are in a new HeaderBlock
					newSubHeadersBlock = newSubHeadersBlock + currLine + "\n"
					lastHeader = parsedLine
				}
			} 
		} 
	}

	// We return a slice of the formed headerblocks that we formed
	return headerBlocks[:headerBlocksSize]
}

// Outputs all of the headers that the player can select and asks them to select a main/sub
// Returns the headers that the player selected
func selectNewConfigLoc(headers []sshmkr_templates.HeaderBlock) (string, string) {
	var mainHeaderIndex int
	var subHeaderIndex int

	var mainHeader string
	var subHeader string

	for currIndex, currHeader := range headers {
		commentStart := strings.Index(currHeader.GetMainHeader(), " ")
		fmt.Printf("%d.) %s\n", currIndex + 1, currHeader.GetMainHeader()[commentStart:])
	}
	fmt.Print("Select a main header: ")
	fmt.Scanln(&mainHeaderIndex)
	mainHeaderIndex = mainHeaderIndex - 1

	if mainHeaderIndex < len(headers) && mainHeaderIndex >= 0 {
		mainHeader = headers[mainHeaderIndex].GetMainHeader()

		fmt.Println("")
		for currIndex, currSubHeader := range headers[mainHeaderIndex].GetSubHeaders() {
			commentStart := strings.Index(currSubHeader, " ")
			fmt.Printf("%d.) %s\n", currIndex + 1, currSubHeader[commentStart:])
		}
		fmt.Print("Select a sub header: ")
		fmt.Scanln(&subHeaderIndex)
		subHeaderIndex = subHeaderIndex - 1

		if subHeaderIndex <  len(headers[mainHeaderIndex].GetSubHeaders()) && subHeaderIndex >= 0 {
			subHeader = headers[mainHeaderIndex].GetSubHeaders()[subHeaderIndex]
		} else {
			fmt.Println("Error! Invalid Choice! Exiting program...")
			os.Exit(1)
		}

	} else {
		fmt.Println("Error! Invalid Choice! Exiting program...")
		os.Exit(1)
	}

	return mainHeader, subHeader
}

// Helper method that checks if the given line is not a comment or empty
func checkIfValid(line string) bool {
	if line != "" && !strings.Contains(line, COMMENT_IND) && !strings.Contains(line, MAIN_HEADER_IND) && !strings.Contains(line, SUB_HEADER_IND) {
		return true
	}
	return false
}

// Helper function that will be used to verify the passed in hostname
// Returns true if we have a match or we pass in nothing (which defaults to the first template)
func checkIfExistingHostname(checkHostname string, verifiedHostname string) bool {
	if verifiedHostname == "*" {
		// This seems to be a little thing regarding the SSH reader I'm using, which is why this test is here
		// Specifically, this is the default value for all hosts?
		return false
	} else if verifiedHostname == checkHostname {
		return true
	} else if len(checkHostname) <= 0 {
		fmt.Println("No hostname passed in, defaulting to first template...")
		return true
	} else {
		return false
	}
}

//// Global Variables

// Constants
const MAIN_HEADER_IND = "####"
const SUB_HEADER_IND = "##"
const COMMENT_IND = "#"

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

	configFile, configFileContents, configFileDecoded := parseConfigFile(configFlagValue)
	configTemplateFile, _, configTemplateFileDecoded := parseConfigFile(fmt.Sprintf("%s_templates", configFlagValue))
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

			template := readSpecificTemplate(*addSource, configTemplateFileDecoded)
			headers := parseConfigHeaders(configFileContents)
			mainHeader, subHeader := selectNewConfigLoc(headers)
			userAddedConfig := interpolateUserInput(template)
			newOutput := addTemplatedConfig(mainHeader, subHeader, userAddedConfig, configFileContents)
			writeToConfigFile(configFlagValue, newOutput)

			fmt.Println("Sucessfully added new config to ssh_config!")
		case "delete":
			deleteCmd.Parse(os.Args[2:])

			newOutput := removeHostConfig(*deleteSource, configFileContents)
			writeToConfigFile(configFlagValue, newOutput)
			
			fmt.Println("Sucessfully removed config from ssh_config!")
		case "copy":
			copyCmd.Parse(os.Args[2:])

			template := readSpecificTemplate(*copySource, configFileDecoded)
			headers := parseConfigHeaders(configFileContents)
			mainHeader, subHeader := selectNewConfigLoc(headers)
			userAddedConfig := interpolateUserInput(template)
			newOutput := addTemplatedConfig(mainHeader, subHeader, userAddedConfig, configFileContents)
			writeToConfigFile(configFlagValue, newOutput)			

			fmt.Println("Sucessfuly copied one config to another config!")
		case "show":
			showCmd.Parse(os.Args[2:])

			getSpecificHostConfig(*showSource, configFileDecoded)
		case "comment":
			commentCmd.Parse(os.Args[2:])

			newOutput, hasCommented := commentHostConfig(*commentSource, configFileContents)
			writeToConfigFile(configFlagValue, newOutput)
			
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
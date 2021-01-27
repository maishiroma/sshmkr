package sshmkr_reader

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"github.com/kevinburke/ssh_config"
	"sshmkr/templates"
)

// Constants
const MAIN_HEADER_IND = "####"
const SUB_HEADER_IND = "##"
const COMMENT_IND = "#"

// Writes out the passed in string into a new file
func WriteToConfigFile(configLoc string, fileContents string) {
	err := ioutil.WriteFile(configLoc, []byte(fileContents), 644)
	if err != nil {
		fmt.Println("Error! The location", configLoc, " cannot be written!")
		os.Exit(1)
	}
}

// Parses the passed config file location to the program
// Returns the open file, the file contents, and the decoded config file
func ParseConfigFile(configLoc string) (*os.File, []byte ,*ssh_config.Config) {
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

// Takes in a parsed ssh config and outputs all of the relevant header comments
// Returns an array of headerBlocks, which are logical groupings of ssh configs
func ParseConfigHeaders(fileContents []byte) []sshmkr_templates.HeaderBlock {
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

// Returns a ConfigTemplate object that contains information on a given template
func ReadSpecificTemplate(hostname string, config_template *ssh_config.Config) sshmkr_templates.ConfigTemplate {
	for _, host := range config_template.Hosts {
		if  CheckIfExistingHostname(hostname, host.Patterns[0].String()) {
			currIndex := 0
			formatted_template_string := "\nHost %s\n"

			// Because we are manually adding the host value in, we need to account for that
			// in the total length of the template
			lenOfTemplate := len(host.Nodes) + 1
			template_kv := make([]ssh_config.KV, lenOfTemplate)
			
			template_kv[currIndex] = ssh_config.KV{Key: "Host", Value: hostname, Comment: ""}
			currIndex = currIndex + 1
			for _, node := range host.Nodes {
				nodeRendered := strings.TrimLeft(node.String(), " ")
				
				if CheckIfValid(nodeRendered) && currIndex < lenOfTemplate {
					// In order to parse the node's key and value, we use indexing
					// Knowing that we are using the first space as the division, we use that as the point of index
					nodeDivider := strings.Index(nodeRendered, " ")
					nodeKey := nodeRendered[0:nodeDivider]
					nodeValue := nodeRendered[nodeDivider+1:]
					
					formatted_template_string = formatted_template_string + "\t" + nodeKey + " %s \n"
					
					// The order of the interpolation in the format template string also correlates
					// to the order of the values in the array
					template_kv[currIndex] = ssh_config.KV{Key: nodeKey, Value: nodeValue, Comment: ""}
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
	fmt.Println("Cannot find template", hostname, "in config_templates file! Typo maybe?")
	os.Exit(-1)

	// Even though it will never reach here, we have to put a return value here
	return sshmkr_templates.ConfigTemplate{}
}

// Helper method that checks if the given line is not a comment or empty
func CheckIfValid(line string) bool {
	if line != "" && !strings.Contains(line, COMMENT_IND) && !strings.Contains(line, MAIN_HEADER_IND) && !strings.Contains(line, SUB_HEADER_IND) {
		return true
	}
	return false
}

// Helper function that will be used to verify the passed in hostname
// Returns true if we have a match or we pass in nothing (which defaults to the first template)
func CheckIfExistingHostname(checkHostname string, verifiedHostname string) bool {
	if verifiedHostname == "*" {
		// This seems to be a little thing regarding the SSH reader I'm using, which is why this test is here
		// Specifically, this is the default value for all hosts?
		return false
	} else if verifiedHostname == checkHostname {
		fmt.Println("Found host config to use for template,", checkHostname, "...")
		fmt.Println("")
		return true
	} else if len(checkHostname) <= 0 {
		fmt.Println("No hostname passed in, defaulting to first template...")
		fmt.Println("")
		return true
	} else {
		return false
	}
}
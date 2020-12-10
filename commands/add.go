package sshmkr_commands

import (
	"strings"
	"sshmkr/reader"
)

// Adds a new host config to a config file
// Returns the new config file contents
func AddTemplatedConfig(mainHeader string, subHeader string, templateString string, fileContents []byte) string {
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

package sshmkr_commands

import (
	"strings"
)

func EditExisingConfig(origHostName string, templateString string, fileContents []byte) string {
	/*
	*	The logic on this script goes by the following:
	*	1. Search for the hostname that we want to edit.
	*	2. Once we find it, we do an inline replacement in the array
	*	3. We exit the script.
	*/

	foundHost := false
	fileContentsArray := strings.Split(string(fileContents), "\n")
	templateArray := strings.Split(templateString, "\n")
	templateArrayLength := len(templateArray)
	templateArrayIndex := 1	// This is 1 because the first value in this split array is an empty line

	for currIndex, currLine := range fileContentsArray {
		if strings.Contains(currLine, origHostName) {
			foundHost = true
		} 
		if foundHost == true {
			if templateArrayIndex < templateArrayLength {
				// Once we find the host name, we perform in line edits to the config
				fileContentsArray[currIndex] = templateArray[templateArrayIndex]
				templateArrayIndex = templateArrayIndex + 1
			} else {
				// We reached the end of the host config so we just exit the loop
				break;
			}
		}
	}

	newContents := strings.Join(fileContentsArray, "\n")
	return newContents
}
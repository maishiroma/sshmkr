package sshmkr_commands

import (
	"strings"
	"fmt"
	"os"
)

// Removes a specified host config from the ssh_config value
// And returns the updated file content
func RemoveHostConfig(hostname string, fileContents []byte) string {
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
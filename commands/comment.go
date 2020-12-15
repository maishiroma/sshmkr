package sshmkr_commands

import (
	"fmt"
	"os"
	"strings"
	"sshmkr/reader"
)

// Comments/Uncomments a specific host config depending if it was already commented or not
// Return the updated file contents and if it did comment it out
func CommentHostConfig(hostname string, fileContents []byte) (string, bool) {
	if len(hostname) <= 0 {
		fmt.Println("Source flag is empty! Please pass in a valid hostname to comment in/out!")
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
		fmt.Println("Cannot find host", hostname, "in config. Typo maybe?")
		os.Exit(-1)
	}

	return strings.Join(fileContentLines, "\n"), hasCommented
}
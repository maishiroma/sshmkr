package sshmkr_input

import (
	"fmt"
	"strings"
	"os"
	"sshmkr/templates"
)

// Takes in a templated string and user input to return a filled host config
func InterpolateUserInput(template sshmkr_templates.ConfigTemplate) (string, string){
	hostName := ""
	templateString := template.GetTemplatedString()
	
	fmt.Println("")
	fmt.Println("~ Template ~")
	for currIndex := 0; currIndex < template.GetNumKeyPairs(); currIndex = currIndex + 1 {
		templateData := template.GetKeyPair(currIndex)
		var userInput string

		fmt.Printf("Enter a value for %s [ default: %s ]: ", templateData.Key, templateData.Value)
		fmt.Scanln(&userInput)
		if userInput == "" {
			templateString = strings.Replace(templateString, "%s", templateData.Value, 1)
			if currIndex == 0 {
				hostName = templateData.Value
			}
		} else {
			templateString = strings.Replace(templateString, "%s", userInput, 1)
			if currIndex == 0 {
				hostName = userInput
			}
		}
	}

	fmt.Println("")
	return templateString, hostName
}

// Outputs all of the headers that the player can select and asks them to select a main/sub
// Returns the headers that the player selected
func SelectNewConfigLoc(headers []sshmkr_templates.HeaderBlock) (string, string) {
	var mainHeaderIndex int
	var subHeaderIndex int

	var mainHeader string
	var subHeader string

	fmt.Println("~ Main Header Selection ~")
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
		fmt.Println("~ Sub Header Selection ~")
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
			fmt.Println("Invalid choice! Exiting program...")
			os.Exit(1)
		}

	} else {
		fmt.Println("Invalid Choice! Exiting program...")
		os.Exit(1)
	}

	return mainHeader, subHeader
}
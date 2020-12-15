package sshmkr_commands

import (
	"fmt"
	"os"
	"sshmkr/reader"
	"github.com/kevinburke/ssh_config"
)

// Prints out a specific host configuration out to standard output
func GetSpecificHostConfig(hostname string, parsedConfig *ssh_config.Config) {
	if len(hostname) <= 0 {
		fmt.Println("Source flag is empty! Please pass in a valid hostname to show!")
		os.Exit(-1)
	}

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
			return
		}
	}

	// We only come here if we cannot find the host specified
	fmt.Println("Cannot find host", hostname, "in config. Typo maybe?")
	os.Exit(-1)
}

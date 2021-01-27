package sshmkr_help

import (
	"flag"
	"fmt"
	"os"
)

func SetHelpContext(f *flag.FlagSet, command string) {
	f.Usage = func() {
		helpText := ""
		switch command {
			case "add":
				helpText = `
Adds in a new SSH configuration to the ssh_config file.

This subcommand starts up an interactive add to the ssh_config file.
The format of the new addition is based of on what is stored in
~/.ssh/config_templates (the default location). One can also omit that
template flag to use the first template structure in that config.

This command also allows for the config to be placed in specific areas of the config, which depend on specific headers. 
These will be pre-determined during runtime and the user will be free to select them.

This command will ignore templates that are commented out.

Example:
  sshmkr add -source nameOfTemplate

Command Flags:
	-source:	Tne name of the source template to use.
 
Global Flags:
	-help: 		Displays the help page for a specific command (or generally)
	-version:	Prints out the current version of the application
	-path:		Changes the default path to look for the ssh_config (default: ~/.ssh/config)
`
			case "delete":
				helpText = `
Deletes a specified SSH Host from the config file.

This command searches from the top to bottom of an ssh_config to determine what to delete. 
This command automatically ignores all hosts that are commented out.

Example:
  sshmkr delete -source nameOfHost

Command Flags:
	-source:	The name of the host to delete (REQUIRED)

Global Flags:
	-help: 		Displays the help page for a specific command (or generally)
	-version:	Prints out the current version of the application
	-path:		Changes the default path to look for the ssh_config (default: ~/.ssh/config)
`
			case "copy":
				helpText = `
Copies an existing SSH host config to use as a new basis for a new one.

This command is useful to duplicate SSH configs that are similar to one
another aside from one or two fields. In the new SSH config, it can even
use the same fields as the original one if needed.

This command, like add, also allows one to specify where to place said 
config in the SSH config file, which is based off on headers.

Example:
  sshmkr copy -source nameOfOriginalHost

Command Flags:
	-source:	The name of the original SSH host to use as a template (REQUIRED)

Global Flags:
	-help: 		Displays the help page for a specific command (or generally)
	-version:	Prints out the current version of the application
	-path:		Changes the default path to look for the ssh_config (default: ~/.ssh/config)
`
			case "show":
				helpText = `
Shows a given SSH config to the screen.

This command outputs the current configuration of a specific SSH host to the screen. 
This is formatted as it appears in the config file as well as to stdout, making it easy to chain into other CLI commands.

Note that if the specified host is commented out, this will ignore said hostname.

Example:
  sshmkr show -source nameOfHost

Command Flags:
	-source: The name of the SSH host to show (REQUIRED)

Global Flags:
	-help: 		Displays the help page for a specific command (or generally)
	-version:	Prints out the current version of the application
	-path:		Changes the default path to look for the ssh_config (default: ~/.ssh/config)
`
			case "comment":
				helpText = `
Comment out a particular SSH host in the SSH config.

This commands searches in the SSH config from top to bottom to either
comment in/out a host configuration. The behavior depends on whether the
host config is already commented in/out.

This is useful for either making specific configs that are relatively similar be active, 
deactivate a particular config, or prevent that config from being parsed in future commands.

Example:
  sshmkr comment -source nameOfHost

Command Flags:
	-source:  The host to comment in/out

Global Flags:
	-help: 		Displays the help page for a specific command (or generally)
	-version:	Prints out the current version of the application
	-path:		Changes the default path to look for the ssh_config (default: ~/.ssh/config)
`
			case "edit":
				helpText = `
Edit one or more of the parameters of an existing SSH config.

When editing an SSH config, the original values will be the default values when prompted,
at which one can either accept them or type in a new value.

Like the other commands, if the host config is commeted out, this command will ignore said
hostname in its search.

Example:
  sshmkr edit -source nameOfHost

Command Flags:
	-source:  The host to comment in/out

Global Flags:
	-help: 		Displays the help page for a specific command (or generally)
	-version:	Prints out the current version of the application
	-path:		Changes the default path to look for the ssh_config (default: ~/.ssh/config)
`
		}
		fmt.Println(helpText)
		os.Exit(0)
	}
}

func DefaultHelp() {
	helpText := `
sshmkr : A Go Binary that respects the ssh_config!

This simple tool aims to help one manage their ssh_config by abiding to
their current ssh_config formatting.

This program (as of now) relies for the ssh_config to be organized by 
the following comments, which the program refers to as headers.
	
	#### Main Header
	## Sub Header
	# Comment

This program also takes a special template file, (defaults to ~/.ssh/config_templates),
that the program can use to create new host configurations from. The format of these are
identical to how a normal host config would look like.

Commands:
	add: 		Adds a new host config into the ssh_config
	delete:		Removes a specified host config from the ssh_config
	comment:	(Un)comments a specified host config from the ssh_config
	copy:		Copies an existing host config and uses it as a template for a new config
	show:		Displays a specified host config
	edit:		Edits an existing SSH config

Global Flags:
	-help: 		Displays the help page for a specific command (or generally)
	-version:	Prints out the current version of the application
	-path:		Changes the default path to look for the ssh_config (default: ~/.ssh/config)
`
	fmt.Println(helpText)
	os.Exit(0)
}

// Prints out the version of the script
func PrintVersion() {
	fmt.Println("sshmkr v1.1.0")
	os.Exit(0)
}
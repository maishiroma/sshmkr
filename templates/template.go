package sshmkr_templates

import (
	"github.com/kevinburke/ssh_config"
)

// Data struct that holds information on all of the headers that exist in the program
type HeaderBlock struct {
	MainHeader string
	SubHeaders []string
}

// Data struct that holds information regarding templated values
type ConfigTemplate struct {
	KeyPairs []ssh_config.KV
	FormattedString string
}

// Returns a templated string that represents the config
func (temp ConfigTemplate) GetTemplatedString() string {
	return temp.FormattedString
}

// Returns a specific key pair from the template
// Key = 0; Value (default) = 1
func (temp ConfigTemplate) GetKeyPair(index int) ssh_config.KV {
	return temp.KeyPairs[index]	
}

// Gets the number of key pairs that are in the template
func (temp ConfigTemplate) GetNumKeyPairs() int {
	return len(temp.KeyPairs)
}

// Gets the main header for that block
func (header HeaderBlock) GetMainHeader() string {
	return header.MainHeader
}

// Gets the sub headers for that block
func (header HeaderBlock) GetSubHeaders() []string {
	return header.SubHeaders
}
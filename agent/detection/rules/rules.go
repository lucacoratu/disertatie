package detection

import (
	"io"

	"gopkg.in/yaml.v2"
)

// Holds all the information in the info field of the rule.yaml files
type YamlInfo struct {
	Name           string `yaml:"name"`           //The name of the rule
	Description    string `yaml:"description"`    //The description of the rule
	Severity       string `yaml:"severity"`       //The severity of the rule, in the string representation
	Classification string `yaml:"classification"` //The classification if it matches, in the string representation
}

// Holds all the modes the search can be made
type YamlSearchMode struct {
	Match string `yaml:"match"` //The string to match exactly
	Regex string `yaml:"regex"` //The regex used for searching
}

type YamlHeaders struct {
	Name  string `yaml:"name"`  //The name of the search to search for matches
	Match string `yaml:"match"` //The string to match exactly
	Regex string `yaml:"regex"` //The regex used for searching
}

// Holds all the information in the request field of the rule YAML file
type YamlRequest struct {
	Method  *YamlSearchMode `yaml:"method"`  //The modes to search on the method
	URL     *YamlSearchMode `yaml:"url"`     //The modes to search on the URL
	Headers []*YamlHeaders  `yaml:"headers"` //The headers to be checked
	Body    *YamlSearchMode `yaml:"body"`    //The string to search for in the body
}

// Holds all the information in the response field of the rule YAML file
type YamlResponse struct {
	Code    *YamlSearchMode `yaml:"code"`    //The modes to search on the status code
	Headers []*YamlHeaders  `yaml:"headers"` //The headers to be checked
	Body    *YamlSearchMode `yaml:"body"`    //The string to search for in the body
}

// Structure which holds all the information about the rule parsed from the rule.yaml file
type YamlRule struct {
	Id       string        `yaml:"id"`       //The ID of the rule (should be unique)
	Info     *YamlInfo     `yaml:"info"`     //The info structure
	Request  *YamlRequest  `yaml:"request"`  //The request matchers
	Response *YamlResponse `yaml:"response"` //The response matchers
}

// Function to read the yaml rule from a reader into the struct
func (yr *YamlRule) FromYAML(r io.Reader) error {
	d := yaml.NewDecoder(r)
	return d.Decode(yr)
}

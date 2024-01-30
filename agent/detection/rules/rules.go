package detection

import (
	"io"

	"gopkg.in/yaml.v2"
)

// Holds all the information in the info field of the rule.yaml files
type RuleInfo struct {
	Name           string `yaml:"name"`           //The name of the rule
	Description    string `yaml:"description"`    //The description of the rule
	Severity       string `yaml:"severity"`       //The severity of the rule, in the string representation
	Classification string `yaml:"classification"` //The classification if it matches, in the string representation
}

// Holds all the modes the search can be made
type RuleSearchMode struct {
	Match string `yaml:"match"` //The string to match exactly
	Regex string `yaml:"regex"` //The regex used for searching
}

// Holds all the information about headers
type HeadersRule struct {
	Name  string `yaml:"name"`  //The name of the search to search for matches
	Match string `yaml:"match"` //The string to match exactly
	Regex string `yaml:"regex"` //The regex used for searching
}

// Holds all the information about request parameters
type RequestParametersRule struct {
	Name  string `yaml:"name"`  //The name of the query variable (can be any which means look through all the query variable names for a match)
	Match string `yaml:"match"` //The string to match exactly
	Regex string `yaml:"regex"` //The regex used for searching
}

// Holds all the information about the body
type BodyRule struct {
	SHA256Sum string `yaml:"sha256sum"` //The SHA256 hash of the body to match
	MD5Sum    string `yaml:"md5sum"`    //The MD5 hash of the body
	Match     string `yaml:"match"`     //The string to match exactly
	Regex     string `yaml:"regex"`     //The regex used for searching
}

// Holds all the information in the request field of the rule YAML file
type RequestRule struct {
	Method     *RuleSearchMode          `yaml:"method"`  //The modes to search on the method
	URL        *RuleSearchMode          `yaml:"url"`     //The modes to search on the URL
	Headers    []*HeadersRule           `yaml:"headers"` //The headers to be checked
	Parameters []*RequestParametersRule `yaml:"params"`  //The request parameters (both from URL and body)
	Body       []*BodyRule              `yaml:"body"`    //The string to search for in the body
}

// Holds all the information in the response field of the rule YAML file
type ResponseRule struct {
	Code    *RuleSearchMode `yaml:"code"`    //The modes to search on the status code
	Headers []*HeadersRule  `yaml:"headers"` //The headers to be checked
	Body    []*BodyRule     `yaml:"body"`    //The string to search for in the body
}

// Structure which holds all the information about the rule parsed from the rule.yaml file
type Rule struct {
	Id       string        `yaml:"id"`       //The ID of the rule (should be unique)
	Info     *RuleInfo     `yaml:"info"`     //The info structure
	Request  *RequestRule  `yaml:"request"`  //The request matchers
	Response *ResponseRule `yaml:"response"` //The response matchers
}

// Function to read the yaml rule from a reader into the struct
func (yr *Rule) FromYAML(r io.Reader) error {
	d := yaml.NewDecoder(r)
	return d.Decode(yr)
}

// Holds all the information about the matched hash of the body
type BodyHashMatch struct {
	BodyHash          string
	BodyHashAlgorithm string
}

package detection

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lucacoratu/disertatie/agent/data"
	"github.com/lucacoratu/disertatie/agent/logging"
)

// Loads all the rules that can be found in the specified directory
// Pass the logger as a parameter for better view of the problems
// @param rulesDirectory - the directory from which the rules should be pulled
// @param logger - the logger to be used to display the errors
// If the directory cannot be opened to read all the files in it then an error is returned
func LoadRulesFromDirectory(rulesDirectory string, logger logging.ILogger) ([]Rule, error) {
	//Check if the directory exists
	_, err := os.Stat(rulesDirectory)
	if err != nil {
		return nil, errors.New("rules directory does not exist")
	}
	//Traverse the directory to get all the rules and append them to the list
	rulesList := make([]Rule, 0)
	err = filepath.WalkDir(rulesDirectory, func(path string, d fs.DirEntry, err error) error {
		//If the directory entry is not a directory
		if !d.IsDir() {
			//If the file ends doesn't .yaml
			if !strings.HasSuffix(d.Name(), ".yaml") {
				logger.Warning("Skipping rule file", path, ", it is not a yaml file, check the file extension")
				return nil
			}
			//The file has .yaml extension
			//Try to load the path content into a rule structure
			rule := Rule{}
			file, err := os.Open(path)
			//Check if an error occured when opening the yaml rule file
			if err != nil {
				//Log the error
				logger.Warning("Skipping rule file", path, "could not open the file for reading", err.Error())
				return nil
			}
			err = rule.FromYAML(file)
			//Check if an error occured when loading the yaml file
			if err != nil {
				logger.Warning("Skipping rule file", path, "error when parsing", err.Error())
				return nil
			}
			//Check if the rule is valid
			err = CheckRule(rule, logger)
			if err != nil {
				//The rule is not valid
				logger.Warning("Skipping rule file", path, "error when checking rule", err.Error())
				return nil
			}
			//Check if the rule id is not already in the list of rules
			var found bool = false
			for _, ruleElem := range rulesList {
				if ruleElem.Id == rule.Id {
					found = true
					break
				}
			}
			if found {
				logger.Warning("Skipping rule file", path, "a rule with this id already exists")
				return nil
			}
			//Add the rule read from file to the list of rules
			rulesList = append(rulesList, rule)
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("could not walk rules directory")
	}
	return rulesList, nil
}

// func PrintRule(rule YamlRule) {

// }

// Check if the rule information is valid or not
// @param info - the rule information structure
// Returns an error if the info field is not valid
func CheckRuleInfo(info *RuleInfo) error {
	//Check if the rule info structure is missing
	if info == nil {
		return errors.New("rule cannot have empty info")
	}

	//Check if the rule severity is one of the allowed (case insensitive) - low, medium, high, critical
	if strings.ToLower(info.Severity) != "low" && strings.ToLower(info.Severity) != "medium" && strings.ToLower(info.Severity) != "high" && strings.ToLower(info.Severity) != "critical" {
		return errors.New("rule severity cannot be something apart from: low, medium, high, critical")
	}

	return nil
}

// Checks if the rule is valid or not
// @param rule - the rule struct to be tested
// @param logger - the logger to be used
// Returns an error which specifies the problem with the rule
func CheckRule(rule Rule, logger logging.ILogger) error {
	//Check if the info field is valid
	if err := CheckRuleInfo(rule.Info); err != nil {
		return err
	}

	//Check if the regexes specified in the rule are compiling
	//Check if the URL regex compiles
	if rule.Request.URL != nil {
		if _, err := regexp.Compile(rule.Request.URL.Regex); err != nil {
			return errors.New("cannot compile URL regex, " + err.Error())
		}
	}
	//Check if the method regex compiles
	if rule.Request.Method != nil {
		if _, err := regexp.Compile(rule.Request.Method.Regex); err != nil {
			return errors.New("cannot compile Method regex, " + err.Error())
		}
	}
	//Check if the parameters regex compiles
	if rule.Request.Parameters != nil {
		for _, parameter := range rule.Request.Parameters {
			if _, err := regexp.Compile(parameter.Regex); err != nil {
				return errors.New("cannot compile regex for parameter, " + parameter.Name + ", " + err.Error())
			}
		}
	}
	//Check if the headers regex compiles
	if rule.Request.Headers != nil {
		for _, header := range rule.Request.Headers {
			if _, err := regexp.Compile(header.Regex); err != nil {
				return errors.New("cannot compile regex for header, " + header.Name + ", " + err.Error())
			}
		}
	}
	//Check if the regex for body
	if rule.Request.Body != nil {
		for _, bodyRule := range rule.Request.Body {
			if _, err := regexp.Compile(bodyRule.Regex); err != nil {
				return errors.New("cannot compile regex for the body, " + err.Error())
			}
		}
	}

	return nil
}

// Converts the severity string to integer value
// @param severity - the string format of the representation (can be low, medium, high, critical) - case insensitive
// Returns the integer coresponding to severity or -1 if the severity type does not exist
func ConvertSeverityStringToInteger(severity string) int64 {
	switch strings.ToLower(severity) {
	case "low":
		return data.LOW
	case "medium":
		return data.MEDIUM
	case "high":
		return data.HIGH
	case "critical":
		return data.CRITICAL
	default:
		return -1
	}
}

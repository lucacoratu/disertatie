package detection

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucacoratu/disertatie/agent/logging"
)

// Loads all the rules that can be found in the specified directory
// Pass the logger as a parameter for better view of the problems
// If the directory cannot be opened to read all the files in it then an error is returned
func LoadRulesFromDirectory(rulesDirectory string, logger logging.ILogger) ([]YamlRule, error) {
	//Check if the directory exists
	_, err := os.Stat(rulesDirectory)
	if err != nil {
		return nil, errors.New("rules directory does not exist")
	}
	//Traverse the directory to get all the rules and append them to the list
	rulesList := make([]YamlRule, 0)
	err = filepath.WalkDir(rulesDirectory, func(path string, d fs.DirEntry, err error) error {
		//If the directory entry is not a directory
		if !d.IsDir() {
			//If the file contains .yaml
			if !strings.HasSuffix(d.Name(), ".yaml") {
				logger.Warning("Skipping rule file", path, ", it is not a yaml file, check the file extension")
				return nil
			}
			//The file has .yaml extension
			//Try to load the path content into a rule structure
			rule := YamlRule{}
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

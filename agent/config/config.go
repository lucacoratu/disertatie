package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/lucacoratu/disertatie/agent/utils"
)

// Structure that will hold the configuration parameters of the proxy
type Configuration struct {
	ListeningProtocol      string `json:"protocol" validate:"required"`                               //The protocol the agent uses to communicate to users
	ListeningAddress       string `json:"address" validate:"required,ipv4"`                           //Address to listen on (127.0.0.1, 0.0.0.0, etc.)
	ListeningPort          string `json:"port" validate:"required,number,gt=0,lt=65536"`              //Port to listen on
	ForbiddenPagePath      string `json:"forbiddenPagePath" validate:"required"`                      //Forbidden page location
	BlacklistUserAgentPath string `json:"blacklistUserAgentPath" validate:"required"`                 //Path to the wordlist of banned User-Agents
	ForwardServerProtocol  string `json:"forwardServerProtocol" validate:"required"`                  //Protocol used when forwarding request to webserver
	ForwardServerAddress   string `json:"forwardServerAddress" validate:"required,ipv4"`              //Address of the webserver to send the request to
	ForwardServerPort      string `json:"forwardServerPort" validate:"required,number,gt=0,lt=65536"` //Port to forward the request to
	APIProtocol            string `json:"apiProtocol"`                                                //API protocol
	APIIpAddress           string `json:"apiIpAddress"`                                               //API ip address
	APIPort                string `json:"apiPort"`                                                    //API port
	RulesDirectory         string `json:"rulesDirectory"`                                             //The directory where rules can be found
	//Mode of operation
}

// Load the configuration from a file
func (conf *Configuration) LoadConfigurationFromFile(filePath string) error {
	//Check if the file exists
	found := utils.CheckFileExists(filePath)
	if !found {
		return errors.New("configuration file cannot be found")
	}
	//Open the file and load the data into the configuration structure
	file, err := os.Open(filePath)
	//Check if an error occured when opening the file
	if err != nil {
		return err
	}
	err = conf.FromJSON(file)
	//Check if an error occured when loading the json from file
	if err != nil {
		return err
	}
	//Initialize the validator of the json data
	validate := validator.New(validator.WithRequiredStructEnabled())
	//Validate the fields of the struct
	err = validate.Struct(conf)
	return err
}

// Convert from json into the configuration structure
func (conf *Configuration) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(conf)
}

// Convert to json the configuration structure
func (conf *Configuration) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(conf)
}

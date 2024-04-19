package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/lucacoratu/disertatie/api/utils"
)

// Structure that will hold the configuration parameters of the proxy
type Configuration struct {
	ListeningAddress    string   `json:"address" validate:"required,ipv4"`              //Address to listen on (127.0.0.1, 0.0.0.0, etc.)
	ListeningPort       string   `json:"port" validate:"required,number,gt=0,lt=65536"` //Port to listen on
	CassandraNodes      []string `json:"cassandraNodes" validate:"required"`            //Cassandra nodes
	CassandraKeyspace   string   `json:"cassandraKeyspace" validate:"required"`         //Cassandra database (keyspace)
	ElasticURL          string   `json:"elasticUrl" validate:"required"`                //Elasticsearch URL (http://<ip>:9200)
	ElasticIndex        string   `json:"elasticIndex" validate:"required"`              //Elastic search index to insert documents into
	ExploitTemplatePath string   `json:"exploitTemplatePath" validate:"required"`       //The template used for generating the exploit
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

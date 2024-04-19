package database

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/data"
	"github.com/lucacoratu/disertatie/api/logging"
)

type ElasticConnection struct {
	configuration config.Configuration
	connection    *elasticsearch.Client
	logger        logging.ILogger
}

func NewElasticConnection(logger logging.ILogger, configuration config.Configuration) *ElasticConnection {
	return &ElasticConnection{logger: logger, configuration: configuration}
}

func (elastic *ElasticConnection) create() error {
	//elastic.connection.Indices.Delete([]string{"logs"})
	//Create the index where the documents will be stored
	_, err := elastic.connection.Indices.Create(elastic.configuration.ElasticIndex)
	//Return the error
	return err
}

func (elastic *ElasticConnection) Init() error {
	//Initialize the config object for the elasticsearch connection
	cfg := elasticsearch.Config{
		Addresses: []string{
			elastic.configuration.ElasticURL,
		},
	}
	//Create the client
	eclient, err := elasticsearch.NewClient(cfg)
	//Check if an error occured when connecting to elasticsearch
	if err != nil {
		return err
	}
	//Save the client object in the elastic connection struct
	elastic.connection = eclient
	//Create the objects needed by the application
	err = elastic.create()
	return err
}

func (elastic *ElasticConnection) InsertLog(log data.LogData) error {
	//Convert request from base64 to string
	rawRequest, err := b64.StdEncoding.DecodeString(log.Request)
	//Check if an error occured when decoding the request from base64
	if err != nil {
		return errors.New("could not decode the request from base64, " + err.Error())
	}

	//Convert response from base64 to string
	rawResponse, err := b64.StdEncoding.DecodeString(log.Response)
	//Check if an error occured when decoding the request from base64
	if err != nil {
		return errors.New("could not decode the request from base64, " + err.Error())
	}

	//Set the decoded values to the log struct
	log.Request = string(rawRequest)
	log.Response = string(rawResponse)

	//Marshal the struct to json string
	logString, err := json.Marshal(log)
	if err != nil {
		return err
	}

	//Insert the document in the elasticsearch index
	_, err = elastic.connection.Index("logs", bytes.NewReader(logString))

	return err
}

type sourceElastic struct {
	Source data.LogData `json:"_source"`
}

type hitsElastic struct {
	Hits []sourceElastic `json:"hits"`
}

type elasticAgentLogs struct {
	Hits hitsElastic `json:"hits"`
}

// Convert json data to LogData structure
func (eal *elasticAgentLogs) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(eal)
}

// Convert LogData structure to json string
func (eal *elasticAgentLogs) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(eal)
}

func (elastic *ElasticConnection) GetRecentLogs() []data.LogDataElastic {
	//Create the query to extract all the logs of the agent
	query := `
	{
		"size": 10,
		"query": { 
			"match_all": {} 
		},
		"sort": [
			{ "timestamp": "desc" }
		]
	}
	`
	//Search the logs in the elasticsearch database
	res, _ := elastic.connection.Search(
		elastic.connection.Search.WithIndex(elastic.configuration.ElasticIndex),
		elastic.connection.Search.WithBody(strings.NewReader(query)),
	)

	//Create the response object
	response := elasticAgentLogs{}
	//Parse the response from json into a struct
	err := response.FromJSON(res.Body)

	//Check if an error occured when parsing the response from json string to struct
	if err != nil {
		return nil
	}

	//Create the return data slice
	returnData := make([]data.LogDataElastic, 0)

	for _, hit := range response.Hits.Hits {
		//Create the request preview
		request_preview := strings.Split(hit.Source.Request, "\n")[0]
		//Create the response preview
		response_preview := strings.Split(hit.Source.Response, "\n")[0]
		//Create the rule findings structure
		returnData = append(returnData, data.LogDataElastic{Id: hit.Source.Id, AgentId: hit.Source.AgentId, RemoteIP: hit.Source.RemoteIP, Timestamp: hit.Source.Timestamp, RequestPreview: request_preview, ResponsePreview: response_preview, Findings: hit.Source.Findings, RuleFindings: hit.Source.RuleFindings})
	}

	elastic.logger.Debug(len(response.Hits.Hits))

	return returnData
}

func (elastic *ElasticConnection) GetLogsPaginated(agentId string) []data.LogDataElastic {
	//Create the query to extract all the logs of the agent
	query := fmt.Sprintf(`
	{
		"size": 10000,
		"query": { 
			"match": {
				"agentId": "%s" 
			} 
		},
		"sort": [
			{ "timestamp": "desc" }
		]
	}
	`, agentId)

	//Search the logs in the elasticsearch database
	res, _ := elastic.connection.Search(
		elastic.connection.Search.WithIndex(elastic.configuration.ElasticIndex),
		elastic.connection.Search.WithBody(strings.NewReader(query)),
	)

	//Create the response object
	response := elasticAgentLogs{}
	//Parse the response from json into a struct
	err := response.FromJSON(res.Body)

	//Check if an error occured when parsing the response from json string to struct
	if err != nil {
		return nil
	}

	//Create the return data slice
	returnData := make([]data.LogDataElastic, 0)

	for _, hit := range response.Hits.Hits {
		//Create the request preview
		request_preview := strings.Split(hit.Source.Request, "\n")[0]
		//Create the response preview
		response_preview := strings.Split(hit.Source.Response, "\n")[0]
		//Create the rule findings structure
		returnData = append(returnData, data.LogDataElastic{Id: hit.Source.Id, AgentId: hit.Source.AgentId, RemoteIP: hit.Source.RemoteIP, Timestamp: hit.Source.Timestamp, RequestPreview: request_preview, ResponsePreview: response_preview, Findings: hit.Source.Findings, RuleFindings: hit.Source.RuleFindings})
	}

	elastic.logger.Debug(len(response.Hits.Hits))

	return returnData
}

type ResponseCountElastic struct {
	Count int64 `json:"count"`
}

func (rce *ResponseCountElastic) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(rce)
}

func (elastic *ElasticConnection) GetTotalCountLogs() (int64, error) {
	client := http.Client{}
	//Create the URL
	url := fmt.Sprintf("%s/%s/_count", elastic.configuration.ElasticURL, elastic.configuration.ElasticIndex)
	response, err := client.Get(url)
	if err != nil {
		return -1, err
	}
	rce := ResponseCountElastic{}
	err = rce.FromJSON(response.Body)
	return rce.Count, err
}

type Bucket struct {
	Key   string `json:"key"`
	Count int64  `json:"doc_count"`
}

type Langs struct {
	Buckets []Bucket `json:"buckets"`
}

type Aggregation struct {
	Langs Langs `json:"langs"`
}

type RuleFindingsAggregationResponse struct {
	Aggregation Aggregation `json:"aggregations"`
}

func (rfar *RuleFindingsAggregationResponse) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(rfar)
}

func (elastic *ElasticConnection) GetRuleFindingsStats() ([]data.FindingsMetrics, error) {
	query := `
	{
		"size": 0,
		"aggs" : {
			"langs" : {
				"terms" : { "field" : "ruleFindings.request.classification.keyword"}
			}
		}}
	`
	//Search the logs in the elasticsearch database
	res, _ := elastic.connection.Search(
		elastic.connection.Search.WithIndex(elastic.configuration.ElasticIndex),
		elastic.connection.Search.WithBody(strings.NewReader(query)),
	)

	response := RuleFindingsAggregationResponse{}
	err := response.FromJSON(res.Body)
	if err != nil {
		return nil, err
	}

	metrics := make([]data.FindingsMetrics, 0)
	for _, metric := range response.Aggregation.Langs.Buckets {
		metrics = append(metrics, data.FindingsMetrics{Classification: metric.Key, Count: metric.Count})
	}

	return metrics, nil
}

func (elastic *ElasticConnection) GetRuleIdStats() ([]data.FindingsMetrics, error) {
	query := `
	{
		"size": 0,
		"aggs" : {
			"langs" : {
				"terms" : { "field" : "ruleFindings.request.ruleId.keyword"}
			}
		}
	}
	`

	//Search the logs in the elasticsearch database
	res, _ := elastic.connection.Search(
		elastic.connection.Search.WithIndex(elastic.configuration.ElasticIndex),
		elastic.connection.Search.WithBody(strings.NewReader(query)),
	)

	response := RuleFindingsAggregationResponse{}
	err := response.FromJSON(res.Body)
	if err != nil {
		return nil, err
	}

	metrics := make([]data.FindingsMetrics, 0)
	for _, metric := range response.Aggregation.Langs.Buckets {
		metrics = append(metrics, data.FindingsMetrics{Classification: metric.Key, Count: metric.Count})
	}

	return metrics, nil
}

// func (elastic *ElasticConnection) GetFindingsStats() ([]data.FindingsMetrics, error) {
// 	query := `
// 	{
// 		"size": 0,
// 		"aggs" : {
// 			"langs" : {
// 				"terms" : { "field" : "ruleFindings.request.ruleId.keyword"}
// 			}
// 		}
// 	}
// 	`

// 	//Search the logs in the elasticsearch database
// 	res, _ := elastic.connection.Search(
// 		elastic.connection.Search.WithIndex(elastic.configuration.ElasticIndex),
// 		elastic.connection.Search.WithBody(strings.NewReader(query)),
// 	)

// 	response := RuleFindingsAggregationResponse{}
// 	err := response.FromJSON(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	metrics := make([]data.FindingsMetrics, 0)
// 	for _, metric := range response.Aggregation.Langs.Buckets {
// 		metrics = append(metrics, data.FindingsMetrics{Classification: metric.Key, Count: metric.Count})
// 	}

// 	return metrics, nil
// }

func (elastic *ElasticConnection) GetRuleFindingsCount() (int64, error) {
	client := http.Client{}
	//Create the URL
	// query := `
	// {
	// 	"query": {
	// 		"term": {"}
	// 	}
	// }
	// `
	url := fmt.Sprintf("%s/%s/_count", elastic.configuration.ElasticURL, elastic.configuration.ElasticIndex)
	response, err := client.Get(url)
	if err != nil {
		return -1, err
	}
	rce := ResponseCountElastic{}
	err = rce.FromJSON(response.Body)
	return rce.Count, err
}

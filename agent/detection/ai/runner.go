package detection

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/lucacoratu/disertatie/agent/config"
	"github.com/lucacoratu/disertatie/agent/data"
	detection "github.com/lucacoratu/disertatie/agent/detection/ai/features"
	"github.com/lucacoratu/disertatie/agent/logging"
)

type AIClassifierRunner struct {
	logger        logging.ILogger
	configuration config.Configuration
}

func NewAIClassifierRunner(logger logging.ILogger, configuration config.Configuration) *AIClassifierRunner {
	return &AIClassifierRunner{logger: logger, configuration: configuration}
}

func (acr *AIClassifierRunner) saveFeaturesInDataset(features data.RequestFeatures, filePath string) {
	//Save the features in a csv file specified in the configuration
	datasetFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		acr.logger.Error("Failed to open dataset file", err.Error())
	}

	//Convert the features to a record for csv file
	var csvRow []string
	//UrlLength|NumberParams|NumberSpecialChars|RatioSpecialChars|NumberRoundBrackets|NumberSquareBrackets|NumberCurlyBrackets|NumberApostrophes|NumberQuotationMarks|NumberDots|NumberSlash|NumberBackslash|NumberComma|NumberColon|NumberSemicolon|NumberMinus|NumberPlus|NumberLessGreater|DistanceDots|DistanceSlash|DistanceBackslash|DistanceComma|DistanceColon|DistanceSemicolon|DistanceMinus|DistancePlus
	csvRow = append(csvRow, strconv.FormatInt(features.UrlLength, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberParams, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberSpecialChars, 10))

	csvRow = append(csvRow, fmt.Sprintf("%f", features.RatioSpecialChars))

	csvRow = append(csvRow, strconv.FormatInt(features.NumberRoundBrackets, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberSquareBrackets, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberCurlyBrackets, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberApostrophes, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberQuotationMarks, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberDots, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberSlash, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberBackslash, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberComma, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberColon, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberSemicolon, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberMinus, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberPlus, 10))
	csvRow = append(csvRow, strconv.FormatInt(features.NumberLessGreater, 10))

	csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceDots))
	csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceSlash))
	csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceBackslash))
	csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceComma))
	csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceColon))
	csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceSemicolon))
	csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceMinus))
	csvRow = append(csvRow, fmt.Sprintf("%f", features.DistancePlus))

	//Create csv writer from dataset file
	w := csv.NewWriter(datasetFile)
	defer w.Flush()

	//Write the features in the file
	w.Write(csvRow)
}

func (acr *AIClassifierRunner) RunAIClassifierOnRequest(r *http.Request) string {
	//Initialize the feature extractor
	featuresExtractor := detection.NewFeaturesExtractor(acr.logger, acr.configuration)

	//Extract the features from the request
	features := featuresExtractor.ExtractFeaturesFromRequest(r)

	//Check if the features should be saved in a dataset
	if acr.configuration.CreateDataset {
		acr.saveFeaturesInDataset(features, acr.configuration.DatasetPath)
	}

	acr.logger.Debug("Features", features)

	if acr.configuration.UseAIClassifier {
		//Select the script to be used based on the model specified in the configuration
		model := fmt.Sprintf("detection/ai/scripts/%s.py", acr.configuration.Classifier)

		//Send the features to the model for prediction
		classification, err := exec.Command("/usr/bin/python", model, features.ToString()).Output()
		if err != nil {
			acr.logger.Error("Failed to run command for prediction using AI model", err.Error())
			return "benign"
		}

		acr.logger.Debug("Classification using AI for request is", string(classification))
		return strings.TrimSpace(string(classification))
	}

	return "benign"
}

func (acr *AIClassifierRunner) RunAIClassifierOnResponse(r *http.Response) {
}

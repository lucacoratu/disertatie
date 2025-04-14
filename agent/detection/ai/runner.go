package detection

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/lucacoratu/disertatie/agent/config"
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

func (vr *AIClassifierRunner) RunAIClassifierOnRequest(r *http.Request) {
	//Initialize the feature extractor
	featuresExtractor := detection.NewFeaturesExtractor(vr.logger, vr.configuration)

	//Extract the features from the request
	features := featuresExtractor.ExtractFeaturesFromRequest(r)

	//Check if the features should be saved in a dataset
	if vr.configuration.CreateDataset == true {
		//Save the features in a csv file specified in the configuration
		datasetFile, err := os.OpenFile(vr.configuration.DatasetPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			vr.logger.Error("Failed to open dataset file", err.Error())
		}

		//Convert the features to a record for csv file
		var csvRow []string
		//UrlLength|NumberParams|NumberSpecialChars|NumberRoundBrackets|NumberSquareBrackets|NumberCurlyBrackets|NumberApostrophes|NumberQuotationMarks|NumberDots|NumberSlash|NumberBackslash|DistanceDots|DistanceSlash|DistanceBackslash
		csvRow = append(csvRow, strconv.FormatInt(features.UrlLength, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberParams, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberSpecialChars, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberRoundBrackets, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberSquareBrackets, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberCurlyBrackets, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberApostrophes, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberQuotationMarks, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberDots, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberSlash, 10))
		csvRow = append(csvRow, strconv.FormatInt(features.NumberBackslash, 10))
		csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceDots))
		csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceSlash))
		csvRow = append(csvRow, fmt.Sprintf("%f", features.DistanceBackslash))

		//Create csv writer from dataset file
		w := csv.NewWriter(datasetFile)
		defer w.Flush()

		//Write the features in the file
		w.Write(csvRow)
	}

	vr.logger.Debug("Features", features)
}

func (vr *AIClassifierRunner) RunAIClassifierOnResponse(r *http.Response) {
}

package detection

import (
	"net/http"
	"strings"

	"github.com/lucacoratu/disertatie/agent/config"
	"github.com/lucacoratu/disertatie/agent/data"
	"github.com/lucacoratu/disertatie/agent/logging"
)

type FeaturesExtractor struct {
	logger        logging.ILogger
	configuration config.Configuration
}

func NewFeaturesExtractor(logger logging.ILogger, configuration config.Configuration) *FeaturesExtractor {
	return &FeaturesExtractor{logger: logger, configuration: configuration}
}

func (featuresExtractor *FeaturesExtractor) ExtractFeaturesFromRequest(request *http.Request) data.RequestFeatures {
	//Initialize the features structure
	var features data.RequestFeatures

	//Extract the features from the request
	//Get the URL of the request
	uri := request.URL.RequestURI()
	features.UrlLength = int64(len(uri))

	//Get the parameters
	urlParams := request.URL.Query()
	//Parse the body parameters
	err := request.ParseForm()
	//Check if an error occured when parsing the request body params
	if err != nil {
		featuresExtractor.logger.Error("Failed to parse request body parameters", err.Error())
	}
	bodyParams := request.PostForm
	features.NumberParams = int64(len(urlParams) + len(bodyParams))

	//Loop through all the query params and body params and extract the features
	//Initialize the features from the params
	features.NumberSpecialChars = 0
	features.NumberRoundBrackets = 0
	features.NumberSquareBrackets = 0
	features.NumberCurlyBrackets = 0
	features.NumberApostrophes = 0
	features.NumberQuotationMarks = 0
	features.NumberDots = 0
	features.NumberSlash = 0
	features.NumberBackslash = 0
	features.NumberComma = 0
	features.NumberColon = 0
	features.NumberSemicolon = 0
	features.NumberMinus = 0
	features.NumberPlus = 0
	features.NumberLessGreater = 0

	paramValuesLen := 0

	//All parameters
	for _, paramValues := range urlParams {
		//The params can also be a list of values
		for _, value := range paramValues {
			//Add the len of the param to total len of parameters
			paramValuesLen += len(value)

			for _, ch := range value {

				//If the character is present in the list of special chars
				if strings.ContainsRune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", ch) {
					features.NumberSpecialChars += 1
				}

				//If the character is ( or )
				if ch == '(' || ch == ')' {
					features.NumberRoundBrackets += 1
				}

				//If the character is [ or ]
				if ch == '[' || ch == ']' {
					features.NumberSquareBrackets += 1
				}

				//If the character is { or }
				if ch == '{' || ch == '}' {
					features.NumberCurlyBrackets += 1
				}

				//If the character is '
				if ch == '\'' {
					features.NumberApostrophes += 1
				}

				//If the character is "
				if ch == '"' {
					features.NumberQuotationMarks += 1
				}

				//If the character is .
				if ch == '.' {
					features.NumberDots += 1
				}

				//If the character is /
				if ch == '/' {
					features.NumberSlash += 1
				}

				//If the character is \
				if ch == '\\' {
					features.NumberBackslash += 1
				}

				//If the character is ,
				if ch == ',' {
					features.NumberComma += 1
				}

				//If the character is :
				if ch == ':' {
					features.NumberColon += 1
				}

				//If the character is ;
				if ch == ';' {
					features.NumberSemicolon += 1
				}

				//If the character is -
				if ch == '-' {
					features.NumberMinus += 1
				}

				//If the character is +
				if ch == '+' {
					features.NumberPlus += 1
				}

				if ch == '<' || ch == '>' {
					features.NumberLessGreater += 1
				}
			}

			//Compute distances
			sumDotDistances := 0
			noDotDistances := 0

			sumSlashDistances := 0
			noSlashDistances := 0

			sumBackslashDistances := 0
			noBackslashDistances := 0

			sumCommaDistances := 0
			noCommaDistances := 0

			sumColonDistances := 0
			noColonDistances := 0

			sumSemicolonDistances := 0
			noSemicolonDistances := 0

			sumMinusDistances := 0
			noMinusDistances := 0

			sumPlusDistances := 0
			noPlusDistances := 0

			for index, ch := range value {
				if ch == '.' {
					part := value[index+1:]
					distance := strings.IndexRune(part, ch)
					//If this is the last char of this type then ignore the distance (it will always be -1)
					if distance == -1 && noDotDistances > 0 {
						continue
					}

					sumDotDistances += distance
					noDotDistances += 1
				}

				if ch == '/' {
					part := value[index+1:]
					distance := strings.IndexRune(part, ch)
					//If this is the last char of this type then ignore the distance (it will always be -1)
					if distance == -1 && noSlashDistances > 0 {
						continue
					}

					sumSlashDistances += distance
					noSlashDistances += 1
				}

				if ch == '\\' {
					part := value[index+1:]
					distance := strings.IndexRune(part, ch)
					//If this is the last char of this type then ignore the distance (it will always be -1)
					if distance == -1 && noBackslashDistances > 0 {
						continue
					}

					sumBackslashDistances += distance
					noBackslashDistances += 1
				}

				if ch == ',' {
					part := value[index+1:]
					distance := strings.IndexRune(part, ch)
					//If this is the last char of this type then ignore the distance (it will always be -1)
					if distance == -1 && noCommaDistances > 0 {
						continue
					}

					sumCommaDistances += distance
					noCommaDistances += 1
				}

				if ch == ':' {
					part := value[index+1:]
					distance := strings.IndexRune(part, ch)
					//If this is the last char of this type then ignore the distance (it will always be -1)
					if distance == -1 && noColonDistances > 0 {
						continue
					}

					sumColonDistances += distance
					noColonDistances += 1
				}

				if ch == ';' {
					part := value[index+1:]
					distance := strings.IndexRune(part, ch)
					//If this is the last char of this type then ignore the distance (it will always be -1)
					if distance == -1 && noSemicolonDistances > 0 {
						continue
					}

					sumSemicolonDistances += distance
					noSemicolonDistances += 1
				}

				if ch == '-' {
					part := value[index+1:]
					distance := strings.IndexRune(part, ch)
					//If this is the last char of this type then ignore the distance (it will always be -1)
					if distance == -1 && noMinusDistances > 0 {
						continue
					}

					sumMinusDistances += distance
					noMinusDistances += 1
				}

				if ch == '+' {
					part := value[index+1:]
					distance := strings.IndexRune(part, ch)
					//If this is the last char of this type then ignore the distance (it will always be -1)
					if distance == -1 && noPlusDistances > 0 {
						continue
					}

					sumPlusDistances += distance
					noPlusDistances += 1
				}
			}

			//Check if the number of distances > 0 (meaning if the character was found the string or not)
			if noDotDistances > 0 {
				features.DistanceDots = float64(sumDotDistances) / float64(noDotDistances)
			} else {
				features.DistanceDots = -1
			}

			if noSlashDistances > 0 {
				features.DistanceSlash = float64(sumSlashDistances) / float64(noSlashDistances)
			} else {
				features.DistanceSlash = -1
			}

			if noBackslashDistances > 0 {
				features.DistanceBackslash = float64(sumBackslashDistances) / float64(noBackslashDistances)
			} else {
				features.DistanceBackslash = -1
			}

			if noCommaDistances > 0 {
				features.DistanceComma = float64(sumCommaDistances) / float64(noCommaDistances)
			} else {
				features.DistanceComma = -1
			}

			if noColonDistances > 0 {
				features.DistanceColon = float64(sumColonDistances) / float64(noColonDistances)
			} else {
				features.DistanceColon = -1
			}

			if noSemicolonDistances > 0 {
				features.DistanceSemicolon = float64(sumSemicolonDistances) / float64(noSemicolonDistances)
			} else {
				features.DistanceSemicolon = -1
			}

			if noMinusDistances > 0 {
				features.DistanceMinus = float64(sumMinusDistances) / float64(noMinusDistances)
			} else {
				features.DistanceMinus = -1
			}

			if noPlusDistances > 0 {
				features.DistancePlus = float64(sumPlusDistances) / float64(noPlusDistances)
			} else {
				features.DistancePlus = -1
			}

			// featuresExtractor.logger.Debug("Dot distances", features.DistanceDots)
		}
	}

	if paramValuesLen > 0 {
		features.RatioSpecialChars = float64(features.NumberSpecialChars) / float64(paramValuesLen)
	} else {
		features.RatioSpecialChars = 0.0
	}

	return features
}

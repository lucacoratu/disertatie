package detection

import (
	"net/http"

	"github.com/lucacoratu/disertatie/agent/config"
	"github.com/lucacoratu/disertatie/agent/logging"
)

type FeaturesExtractor struct {
	logger        logging.ILogger
	configuration config.Configuration
}

func NewFeaturesExtractor(logger logging.ILogger, configuration config.Configuration) *FeaturesExtractor {
	return &FeaturesExtractor{logger: logger, configuration: configuration}
}

func (featuresExtractor *FeaturesExtractor) ExtractFeaturesFromRequest(request *http.Request) {

}

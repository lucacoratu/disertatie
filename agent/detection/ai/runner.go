package detection

import (
	"net/http"

	"github.com/lucacoratu/disertatie/agent/logging"
)

type AIClassifierRunner struct {
	logger logging.ILogger
}

func NewAIClassifierRunner(logger logging.ILogger) *AIClassifierRunner {
	return &AIClassifierRunner{logger: logger}
}

func (vr *AIClassifierRunner) RunAIClassifierOnRequest(r *http.Request) {

}

func (vr *AIClassifierRunner) RunAIClassifierOnResponse(r *http.Response) {
}

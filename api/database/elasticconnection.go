package database

import "github.com/lucacoratu/disertatie/api/data"

type IElasticConnection interface {
	Init() error
	InsertLog(log data.LogData) error
	GetLogsPaginated(agentId string) []data.LogDataElastic
	GetRecentLogs() []data.LogDataElastic
	GetTotalCountLogs() (int64, error)
	GetRuleFindingsStats() ([]data.FindingsMetrics, error)
	GetRuleIdStats() ([]data.FindingsMetrics, error)
}

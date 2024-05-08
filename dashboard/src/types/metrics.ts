/*
	Id     int64  `json:"id"`     //The ID of the method metrics
	Method string `json:"method"` //The method name
	Count  int64  `json:"count"`  //The number of occurences of specific method
*/

type MethodsMetrics = {
    id: number,
    method: string,
    count: number
}

/*
 	Id    int64  `json:"id"`    //The ID of the metric
	Date  string `json:"date"`  //The date when the counts have been computed
	Count int64  `json:"count"` //The number of requests till the specified date
 */

type DaysMetrics = {
    id: number,
    date: string,
    count: number,
}

/*
	Id         int64  `json:"id"`         //The ID of the metric
	StatusCode string `json:"statusCode"` //The status code of the metric
	Count      int64  `json:"count"`      //The number of occurences of the status code
 */
type StatusCodeMetrics = {
    id: number,
    statusCode: string,
    count: number,
}

/*
	Id        int64  `json:"id"`        //The ID of the metric
	IPAddress string `json:"ipAddress"` //The IP address of the metric
	Count     int64  `json:"count"`     //The number of occurences of the IP address
*/

type IPMetrics = {
    id: number,
    ipAddress: string,
    count: number,
}

/*
type FindingsCountMetrics struct {
	FindingsCount     int64 `json:"findingsCount"`
	RuleFindingsCount int64 `json:"ruleFindingsCount"`
}
*/

type FindingsCountMetrics = {
    findingsCount: number,
    ruleFindingsCount: number
}

type FindingsMetrics = {
    classification: string,
    count: number
}

type MethodMetricsResponse = {
    metrics: MethodsMetrics[]
}

type DaysMetricsResponse = {
    metrics: DaysMetrics[]
}

type StatusCodeMetricsResponse = {
    metrics: StatusCodeMetrics[],
}

type IPAddressesMetricsResponse = {
    metrics: IPMetrics[],
}

type FindingsMetricsResponse = {
    metrics: FindingsMetrics[],
}

type FindingsCountMetricsResponse = {
    metrics: FindingsCountMetrics
}
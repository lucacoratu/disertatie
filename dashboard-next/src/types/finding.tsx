/*
type FindingClassificationString struct {
	IntegerFormat int64  `json:"intFormat"`    //The integer format
	StringFormat  string `json:"stringFormat"` //The string format
}
*/

/*
type FindingClassificationStringResponse struct {
	FindingsString []data.FindingClassificationString `json:"findingsString"`
}
 */

type FindingClassificationString = {
    intFormat: number,
    stringFormat: string,
    description: string,
}

type FindingClassificationStringResponse = {
    findingsString: FindingClassificationString[],
}
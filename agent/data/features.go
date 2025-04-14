package data

type RequestFeatures struct {
	UrlLength            int64 //The length of the url
	NumberParams         int64 //The number of parameters including the ones from url and from body
	NumberSpecialChars   int64 //The number of special characters from the parameters
	NumberRoundBrackets  int64 //The number of round brackets from the parameters
	NumberSquareBrackets int64 //The number of square brackets from the parameters
	NumberCurlyBrackets  int64 //The number of curly brackets from the parameters
	NumberApostrophes    int64 //The number of apostrophes from the parameters
	NumberQuotationMarks int64 //The number of quotation marks from the parameters
	NumberDots           int64 //The number of dots (.) from the parameters
	NumberSlash          int64 //The number of / from the parameters
	NumberBackslash      int64 //The number of \ from the parameters

	DistanceDots      float64 //Avg number of characters between succesive dots (.) from the parameters
	DistanceSlash     float64 //Avg number of characters between succesive slashes (/) from the parameters
	DistanceBackslash float64 //Avg number of characters between succesive backslashes (\) from the parameters
}

package handlers

import (
	"net/http"

	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/data"
	"github.com/lucacoratu/disertatie/api/database"
	"github.com/lucacoratu/disertatie/api/jwt"
	"github.com/lucacoratu/disertatie/api/logging"

	request "github.com/lucacoratu/disertatie/api/data/request"
	response "github.com/lucacoratu/disertatie/api/data/response"
)

type AuthHandler struct {
	logger            logging.ILogger
	configuration     config.Configuration
	dbConnection      database.IConnection
	elasticConnection database.IElasticConnection
}

// Creates a new handler that will hold the functions necessary for registering proxies
func NewAuthHandler(logger logging.ILogger, configuration config.Configuration, dbConnection database.IConnection, elasticConnection database.IElasticConnection) *AuthHandler {
	return &AuthHandler{logger: logger, configuration: configuration, dbConnection: dbConnection, elasticConnection: elasticConnection}
}

// Function to handle login
func (ah *AuthHandler) Login(rw http.ResponseWriter, r *http.Request) {
	//Check if the method is POST
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Get the credentials from the request body
	creds := request.LoginRequest{}
	err := creds.FromJSON(r.Body)
	//Check if an error occured when parsing the login request body
	if err != nil {
		ah.logger.Error("Could not parse JSON data from body in the coresponding structure", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.PARSE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}

	//Check the provided credentials
	exists, id, err := ah.dbConnection.CheckUserCredentials(creds.Username, creds.Password)
	//Check if an error occured
	if err != nil {
		ah.logger.Error("Could not check the user credentials", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}

	//Check if the credentials are invalid
	if !exists {
		rw.WriteHeader(http.StatusUnauthorized)
		retErr := data.APIError{Code: data.AUTH_ERROR, Message: "Invalid credentials"}
		retErr.ToJSON(rw)
		return
	}

	//Create the jwt and send it back as a cookie
	jwtToken, err := jwt.GenerateJWT(id, creds.Username)
	if err != nil {
		ah.logger.Error("Could not generate JWT token for", creds.Username, err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.JWT_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}

	//Return the jwt token to the client
	response := response.LoginResponse{Token: jwtToken}
	rw.Header().Add("Set-Cookie", "session="+jwtToken)
	rw.WriteHeader(http.StatusOK)
	//Add the cookie for the session so that the browser can set it

	response.ToJSON(rw)
}

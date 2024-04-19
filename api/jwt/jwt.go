package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Key that will be used to sign the jwt tokens
// TO DO... move it to environment variable (NOT SECURE)
var jwtKey = []byte("apisupersecretjwtkey")

/*
 * This structure holds the information stored inside the jwt payload
 */
type APIClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

/*
 * This function will generate a JWT token that will be valid for 60 minutes
 * If an error occurs during the JWT generation then this function will return an error
 * Else the function will return nil for error
 */
func GenerateJWT(id string, username string) (string, error) {
	//Set the expiration time to be in an hour from the time of generation
	expirationTime := time.Now().Add(time.Hour * 1)
	//Create the data of the JWT
	claims := &APIClaims{
		Id:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	//Set the HMAC_SHA256 algorithm to be used with a secret password
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//Sign the token
	return token.SignedString(jwtKey)
}

/*
 * This function will validate the JWT token
 */
func ValidateJWT(signedToken string) (APIClaims, error) {
	claims := APIClaims{}
	token, err := jwt.ParseWithClaims(signedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return APIClaims{}, err
	}

	if !token.Valid {
		return APIClaims{}, errors.New("JWT token is invalid")
	}

	return claims, nil
}

// func ExtractClaims(signedToken string) (*APIClaims, error) {
// 	//Try to parse the token received (general format)
// 	token, err := jwt.ParseWithClaims(
// 		signedToken,
// 		&JWTClaim{},
// 		func(t *jwt.Token) (interface{}, error) {
// 			return []byte(jwtKey), nil
// 		},
// 	)
// 	//Check if an erorr occured during parsing
// 	if err != nil {
// 		//An error occured (invalid format of the JWT)
// 		return nil, err
// 	}

// 	//Try to convert the data in the JWT to the JWTClaim structure
// 	claims, ok := token.Claims.(*JWTClaim)
// 	//Check if the convertion was ok
// 	if !ok {
// 		//The conversion was no ok
// 		return nil, errors.New("couldn't parse the claims")
// 	}

// 	return claims, nil
// }

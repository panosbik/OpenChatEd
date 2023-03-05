package constants

import "errors"

var (
	// General
	ErrServer    = errors.New("we're sorry, but we could not complete your request at this time, please try again later or contact customer support if the issue persists")
	ErrBodyParse = errors.New("an error occurred while processing the body of the application, please check that you have provided valid fields")

	// User
	ErrAccessDenied            = errors.New("access denied, invalid or missing authorization header, please include a valid 'Bearer' token in the 'Authorization' header to access this resource")
	ErrInvalidToken            = errors.New("the validation of your request failed due to an invalid JSON Web Token (JWT), please ensure that the provided token is valid and has not been tampered with")
	ErrInvalidRefreshToken     = errors.New("the validation of your request failed due to an invalid Refresh Token, please ensure that the provided token is valid and has not been tampered with")
	ErrInvalidLoginCredentials = errors.New("we could not log you in at this time, please check that you have provided the correct login credentials and try again")
	ErrEmailVerificationError  = errors.New("email address not verified, please check your email inbox and follow the verification link to confirm your email address")
	ErrDuplicateUsername       = errors.New("the username you provided is already associated with an existing user account")
	ErrDuplicateEmailUser      = errors.New("the email address you provided is already associated with an existing user account")
	ErrInvalidTermSearchUser   = errors.New("your request is missing the required information, check that you have provided the 'term' parameter in your query and try again")
)

package auth

// enum for all possible errors as strings
const (
	ErrInvalidCredentials  = "InvalidCredentials"
	ErrUserNotFound        = "UserNotFound"
	ErrUserExists          = "UserAlreadyExists"
	ErrInvalidPassword     = "InvalidPassword"
	ErrInvalidProvider     = "InvalidProvider"
	ErrMissingCode         = "MissingCode"
	ErrTokenExchangeFailed = "FailedToExchangeToken"
	ErrFailedToGetUser     = "FailedGetUser"
	ErrCreationFailed      = "FailedCreateUser"
	ErrJWTSigningFailed    = "FailedSignJWT"
	ErrJWTVerification     = "FailedVerifyJWT"
	ErrJWTInvalid          = "InvalidJWT"
	ErrJWTExpired          = "ExpiredJWT"
)

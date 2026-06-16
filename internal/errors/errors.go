package errors

import (
	"errors"
	"net/http"

	"github.com/joomcode/errorx"
)

type ErrorType struct {
	StatusCode int
	Type       *errorx.Type
}

var (
	DBNS          = errorx.NewNamespace("DB")
	DBFetchFailed = errorx.NewType(DBNS, "DB_FETCH_FAILED")

	RequestNS      = errorx.NewNamespace("REQUEST")
	InvalidRequest = errorx.NewType(RequestNS, "INVALID_REQUEST")

	ErrBadJSON           = errors.New("bad json")
	ErrInvalidActionType = errors.New("invalid action type")
)

var Error = []ErrorType{
	{Type: DBFetchFailed, StatusCode: http.StatusInternalServerError},
	{Type: InvalidRequest, StatusCode: http.StatusBadRequest},
	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrInvalidUserInput,
	},
	{
		StatusCode: http.StatusForbidden,
		Type:       ErrAccessError,
	},
	{
		StatusCode: http.StatusInternalServerError,
		Type:       ErrInternalServerError,
	},
	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrAuthClient,
	},

	{
		StatusCode: http.StatusUnauthorized,
		Type:       ErrInvalidAccessToken,
	},

	{
		StatusCode: http.StatusUnauthorized,
		Type:       ErrUnAuthorizedAccess,
	},
	{
		StatusCode: http.StatusInternalServerError,
		Type:       ErrUnableToGet,
	},

	{
		StatusCode: http.StatusInternalServerError,
		Type:       ErrUnableTocreate,
	},
	{
		StatusCode: http.StatusNotFound,
		Type:       ErrResourceNotFound,
	},

	{
		StatusCode: http.StatusConflict,
		Type:       ErrDataAlredyExist,
	},
	{
		StatusCode: http.StatusNotFound,
		Type:       ErrNoRecordFound,
	},
	{
		StatusCode: http.StatusInternalServerError,
		Type:       ErrUnableToDeletelError,
	},
	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrProgramStatus,
	},
	{
		StatusCode: http.StatusInternalServerError,
		Type:       ERR_SESSION_CREATION_FAILED,
	}, {
		StatusCode: http.StatusBadRequest,
		Type:       ERR_INVALID_PAYLOAD,
	},

	{
		StatusCode: http.StatusForbidden,
		Type:       ERR_REDIS,
	},

	{
		StatusCode: http.StatusBadRequest,
		Type:       ERR_BAD_REQUEST,
	},

	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrPoolNotStarted,
	},
	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrSessionNotFound,
	},
	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrUnmarshaling,
	},
	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrMarshaling,
	},
	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrSavingSession,
	},
}

var (
	databaseError                = errorx.NewNamespace("database error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	invalidInput                 = errorx.NewNamespace("validation error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	resourceNotFound             = errorx.NewNamespace("not found").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	unauthorized                 = errorx.NewNamespace("unauthorized").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	ineligible                   = errorx.NewNamespace("ineligible").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	AccessDenied                 = errorx.RegisterTrait("You are not authorized to perform the action")
	Ineligible                   = errorx.RegisterTrait("You are not eligible to perform the action")
	serverError                  = errorx.NewNamespace("INTERNAL_SERVIER_ERROR")
	authoriztionClientError      = errorx.NewNamespace("authorization client error")
	Unauthenticated              = errorx.NewNamespace("user authentication failed")
	ProgramError                 = errorx.NewNamespace("program error")
	httpError                    = errorx.NewNamespace("http error")
	dbError                      = errorx.NewNamespace("db error")
	logicFailed                  = errorx.NewNamespace("business logic Failed").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	requestFailed                = errorx.NewNamespace("request binding Failed").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	bodyreadFailed               = errorx.NewNamespace("reading response body Failed").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	unroutableLocation           = errorx.NewNamespace("unroutable location").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	pgtypeJsonbParseError        = errorx.NewNamespace("failed to parse message data")
	statusResponseError          = errorx.NewNamespace("status response error")
	sessionCreationFailed        = errorx.NewNamespace("SESSION_CREATION_FAILED")
	sessionAlreadyExists         = errorx.NewNamespace("SESSION_ALREADY_EXISTS")
	invalidPayload               = errorx.NewNamespace("INVALID_PAYLOAD")
	rediserror                   = errorx.NewNamespace("REDIS_ERROR")
	invalidtoken                 = errorx.NewNamespace("INVALID_TOKEN")
	poolnotstarted               = errorx.NewNamespace("POOL_NOT_STARTED")
	sessionnotfound              = errorx.NewNamespace("SESSION_NOT_FOUND")
	errunmarshaling              = errorx.NewNamespace("ERR_UNMARSHALING")
	errmarshaling                = errorx.NewNamespace("ERR_MARSHALING")
	errsavingsession             = errorx.NewNamespace("ERR_SAVING_Session")
	errdesiredamount             = errorx.NewNamespace("ERR_DESIRED_AMOUNT")
	errestablishment             = errorx.NewNamespace("ERR_Establishment")
	errinvalidpin                = errorx.NewNamespace("ERR_INVALID_PIN")
	erraccountlocked             = errorx.NewNamespace("ERR_ACCOUNT_LOCKED")
	erraccountstatusupdatefailed = errorx.NewNamespace("ERR_ACCOUNT_STATUS_UPDATE_FAILED")
	apicallfailed                = errorx.NewNamespace("ERR_BANK_API_CALL_FAILED")
)

var (
	ErrUnAuthorizedAccess       = errorx.NewType(unauthorized, "unauthorized access")
	ErrFailedToParseMessageData = errorx.NewType(pgtypeJsonbParseError, "failed to parse message data")
	ErrUnableTocreate           = errorx.NewType(databaseError, "unable to create")
	ErrDataAlredyExist          = errorx.NewType(databaseError, "data already exist")
	ErrUnableToGet              = errorx.NewType(databaseError, "unable to get")
	ErrInvalidUserInput         = errorx.NewType(invalidInput, "invalid user input")
	ErrInactiveUserStatus       = errorx.NewType(invalidInput, "Inactive user status")
	ErrResourceNotFound         = errorx.NewType(resourceNotFound, "resource not found")
	ErrAccessError              = errorx.NewType(unauthorized, "Unauthorized", AccessDenied)
	ErrIneligibleError          = errorx.NewType(ineligible, "Ineligible", Ineligible)
	ErrInternalServerError      = errorx.NewType(serverError, "internal server error")
	ErrAuthClient               = errorx.NewType(authoriztionClientError, "authorization client error")
	ErrSSOAuthenticationFailed  = errorx.NewType(Unauthenticated, "user authentication failed")
	ErrInvalidAccessToken       = errorx.NewType(Unauthenticated, "invalid token").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	ErrSSOError                 = errorx.NewType(serverError, "sso communication failed")
	ErrAccountingError          = errorx.NewType(serverError, "accounting error")
	ErrUnExpectedError          = errorx.NewType(serverError, "unexpected error occurred")
	ErrUnableToUpdate           = errorx.NewType(databaseError, "unable to update")
	ErrUnableToDeletelError     = errorx.NewType(databaseError, "could not delete record")
	ErrNoRecordFound            = errorx.NewType(resourceNotFound, "no record found")
	ErrProgramStatus            = errorx.NewType(ProgramError, "program status error")
	ErrProgramAmount            = errorx.NewType(ProgramError, "spending limit error")
	ErrSMSSend                  = errorx.NewType(serverError, "couldn't send sms")
	ErrHTTPRequestPrepareFailed = errorx.NewType(httpError, "couldn't prepare http request")
	ErrWriteError               = errorx.NewType(dbError, "could not write to db")
	ErrReadError                = errorx.NewType(dbError, "could not read data from db")
	ErrLogicFailed              = errorx.NewType(logicFailed, "logic failure")
	ErrHTTPRequestBinding       = errorx.NewType(requestFailed, "binding failure")
	ErrReadingResponseBody      = errorx.NewType(bodyreadFailed, "reading body failure")
	ErrUnroutableLocaiton       = errorx.NewType(unroutableLocation, "Failed to find route for location")
	ErrStatusResponse           = errorx.NewType(statusResponseError, "status response error")
	ErrNoActiveSession          = errorx.NewType(resourceNotFound, "no active session found")
	ErrActiveSessionExists      = errorx.NewType(databaseError, "active session already exists")
	ErrRetrivingSession         = errorx.NewType(resourceNotFound, "error retriving session")
	ErrSetingSession            = errorx.NewType(databaseError, "error setting session")
	ErrClearingSession          = errorx.NewType(databaseError, "error clearing session")

	ERR_BAD_REQUEST             = errorx.NewType(unauthorized, "bad request")
	ERR_TIMESTAMP_INVALID       = errorx.NewType(unauthorized, "invalid timestamp")
	ERR_SESSION_CREATION_FAILED = errorx.NewType(sessionCreationFailed, "create redis session failed")
	ERR_SESSION_ALREADY_EXISTS  = errorx.NewType(sessionAlreadyExists, "session already exists")
	ERR_INVALID_PAYLOAD         = errorx.NewType(invalidPayload, "invalid payload")

	ERR_REDIS = errorx.NewType(rediserror, "Redis failure")

	ERR_INVALID_TOKEN = errorx.NewType(invalidtoken, "invalid token")
	ErrPoolNotStarted = errorx.NewType(poolnotstarted, "pool not started")

	ERR_UNAUTHORIZED   = errorx.NewType(unauthorized, "unauthorized")
	ErrSessionNotFound = errorx.NewType(sessionnotfound, "Session Not Found.")
	ErrUnmarshaling    = errorx.NewType(errunmarshaling, "Err unmarshaling data.")
	ErrMarshaling      = errorx.NewType(errmarshaling, "Error marshaling data.")
	ErrSavingSession   = errorx.NewType(errsavingsession, "Error saving session data.")
)

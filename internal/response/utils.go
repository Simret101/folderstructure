package response

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"folderstructure/internal/errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-playground/validator/v10"
	"github.com/joomcode/errorx"
	"github.com/spf13/viper"
)

func SendSuccessResponse(ctx *gin.Context, statusCode int, code Code, data any, message any) {
	requestID := ctx.Value("X-Request-Id")
	if requestID == "" {
		requestID = "UNKNOWN"
	}

	response := SuccessResponse{
		Status:  "success",
		Code:    string(code),
		Message: message,
		Data:    data,
		Metadata: Metadata{
			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			RequestID:       fmt.Sprintf("%v", requestID),
		},
	}

	ctx.JSON(
		statusCode,
		response,
	)
}

func SendErrorResponse(ctx *gin.Context, err *ErrorResponse) {
	ctx.AbortWithStatusJSON(err.Code, err)
}

func SendErrorResponseFormated(ctx *gin.Context, err error) {
	requestID, ok := ctx.Get("X-Request-Id")
	if !ok {
		requestID = "UNKNOWN"
	}

	statusCode := http.StatusInternalServerError
	message := "Internal Server Error"
	code := "internal_error"
	if err != nil {
		if validationErrors, ok := err.(validation.ValidationErrors); ok {
			statusCode = http.StatusBadRequest
			message = "Validation error"
			code = "validation_error"
			for _, fieldError := range validationErrors {
				message = fieldError.Field() + ": " + fieldError.ActualTag()
			}
		} else {
			for _, e := range errors.Error {
				if errorx.IsOfType(err, e.Type) {
					statusCode = e.StatusCode
					message = strings.ReplaceAll(e.Type.FullName(), fmt.Sprintf("%s.", e.Type.Namespace().String()), "")

					if strings.Contains(err.Error(), "attempt left") || strings.Contains(err.Error(), "attempts left") {
						message = strings.ReplaceAll(err.Error(), "CUSTOMER_NOT_VERIFIED.invalid pin", "Invalid Pin")
					}

					code = e.Type.Namespace().String()
					break
				}
			}

			if message == "Internal Server Error" {
				message = err.Error()
				code = "unknown_error"
			}
		}
	}

	response := ErrorResponseFormat{
		Status:  "error",
		Code:    code,
		Message: message,
		Metadata: Metadata{
			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			RequestID:       fmt.Sprintf("%v", requestID),
		},
	}

	ctx.AbortWithStatusJSON(statusCode, response)
}

func GetErrorFrom(err error) *ErrorResponse {
	debugMode := viper.GetBool("debug")

	for _, e := range errors.Error {
		if errorx.IsOfType(err, e.Type) {
			er := errorx.Cast(err)
			res := ErrorResponse{
				Code:       e.StatusCode,
				Message:    er.Message(),
				FieldError: ErrorFields(er.Cause()),
			}

			if debugMode {
				res.Description = fmt.Sprintf("Error: %v", er)
				res.StackTrace = fmt.Sprintf("%+v", errorx.EnsureStackTrace(err))
			}

			return &res
		}
	}

	return &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Unknown server error",
	}
}

func ErrorFields(err error) []FieldError {
	var errs []FieldError

	if data, ok := err.(validation.ValidationErrors); ok {
		for i, v := range data {
			errs = append(errs, FieldError{
				Name:        fmt.Sprintf("%d", i+1),
				Description: v.Error(),
			},
			)
		}

		return errs
	}

	return nil
}

func SendAuthzResponseErr(ctx *gin.Context, err error, message any) {
	requestID, ok := ctx.Get("X-Request-Id")
	if !ok {
		requestID = "UNKNOWN"
	}

	statusCode := http.StatusInternalServerError
	code := "internal_error"

	if err != nil {
		for _, e := range errors.Error {
			if errorx.IsOfType(err, e.Type) {
				statusCode = e.StatusCode
				code = e.Type.Namespace().String()
				break
			}
		}

		if message == "Internal Server Error" {
			message = err.Error()
			code = "unknown_error"
		}

	}

	response := ErrorAuthzResponseFormat{
		Status:  "error",
		Code:    code,
		Message: message,
		Metadata: Metadata{
			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			RequestID:       fmt.Sprintf("%v", requestID),
		},
	}

	ctx.AbortWithStatusJSON(statusCode, response)
}

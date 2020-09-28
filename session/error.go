package session

import (
	"encoding/json"
	"github.com/astaxie/beego/context"
	"log"
	"net/http"
)

type Error struct {
	Status      int    `json:"status"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func (sessionError Error) Error() string {
	str, err := json.Marshal(sessionError)
	if err != nil {
		log.Panicln(err)
	}
	return string(str)
}

func HandleError(ctx *context.Context, err *Error) {
	ctx.Output.JSON(err, false, false)
}

//func ParseError(err string) (Error, bool) {
//	var sessionErr Error
//	json.Unmarshal([]byte(err), &sessionErr)
//	return sessionErr, sessionErr.Code > 0 && sessionErr.Description != ""
//}

func BadRequestError() *Error {
	description := "The request body can’t be pasred as valid data."
	return createError(http.StatusAccepted, http.StatusBadRequest, description)
}

func ServerError() *Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(http.StatusInternalServerError, http.StatusInternalServerError, description)
}

func NotFoundError() *Error {
	description := "The endpoint is not found."
	return createError(http.StatusAccepted, http.StatusNotFound, description)
}

func AuthorizationError() *Error {
	description := "Unauthorized, maybe invalid token."
	return createError(http.StatusAccepted, 401, description)
}

func ForbiddenError() *Error {
	description := http.StatusText(http.StatusForbidden)
	return createError(http.StatusAccepted, http.StatusForbidden, description)
}

func TransactionError() *Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(http.StatusInternalServerError, 10001, description)
}

func BadDataError() *Error {
	description := "The request data has invalid field."
	return createError(http.StatusAccepted, 10002, description)
}

func AssetForbiddenError() *Error {
	description := "Asset access forbidden."
	return createError(http.StatusAccepted, 10003, description)
}

func InsufficientAccountBalanceError() *Error {
	description := "Insufficient balance."
	return createError(http.StatusAccepted, 20117, description)
}

func BlazeServerError() *Error {
	description := "Blaze server error."
	return createError(http.StatusInternalServerError, 7000, description)
}

func BlazeTimeoutError() *Error {
	description := "The blaze operation timeout."
	return createError(http.StatusInternalServerError, 7001, description)
}

func createError(status, code int, description string) *Error {
	return &Error{
		Status:      status,
		Code:        code,
		Description: description,
	}
}

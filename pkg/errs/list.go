package errs

import (
	"fmt"

	"net/http"
	"strconv"
)

type CodeError struct {
	Code              int
	Message           string
	AdditionalMessage string
	HttpCode          int
}

// newCodeError create new Code error and fill ErrorsDescription
func newCodeError(code int, description, message string, httpCode int) *CodeError {
	if _, ok := ErrorsDescription[code]; ok {
		panic(fmt.Sprintf("error Code duplicate %v", code))
	}
	ErrorsDescription[code] = struct{}{}
	return &CodeError{Code: code, Message: description, HttpCode: httpCode}
}

func (c CodeError) Error() string {
	return strconv.Itoa(c.Code)
}

// AddInfo create new CodeError with AdditionalMessage as s
func (c CodeError) AddInfo(s string) CodeError {
	if c.AdditionalMessage != "" {
		c.AdditionalMessage += "\n\r"
	}
	c.AdditionalMessage += s
	return c
}

func (c CodeError) AddInfoErrMessage(err error) CodeError {
	if c.AdditionalMessage != "" {
		c.AdditionalMessage += "\n\r"
	}
	codeError, ok := err.(CodeError)
	if ok {
		c.AdditionalMessage += codeError.AdditionalMessage
		if codeError.Code != c.Code {
			c.AdditionalMessage += "\n\r"
			c.AdditionalMessage += "Errors cascade!" + codeError.Error()
		}
		return c
	}
	c.AdditionalMessage += err.Error()
	return c
}

// Equal compare code of c and err
func (c CodeError) Equal(err error) bool {
	if err == nil {
		return false
	}
	err = Cause(err)
	codeError, ok := err.(CodeError)
	if !ok {
		return false
	}
	return codeError.Code == c.Code
}

// ErrorsDescription description of every error in bindings
// Fill by newCodeError
var ErrorsDescription = map[int]struct{}{}

var (
	UnknownError           = newCodeError(1, "Что-то пошло не так. Попробуйте позже.", "UnknownError", http.StatusInternalServerError)   // Неизвестная ошибка
	TransportError         = newCodeError(2, "Что-то пошло не так. Попробуйте позже.", "TransportError", http.StatusInternalServerError) // Ошибка на уровне HTTP
	RequestValidationError = newCodeError(100001, "Некорректно указаны данные.", "RequestValidationError", http.StatusBadRequest)        // Запрос не проходит базовую валидацию по полям
	ForbiddenOperation     = newCodeError(403, "Вам не разрешено выполнять это действие", "ForbiddenOperation", http.StatusForbidden)

	InvalidToken               = newCodeError(200001, "Некорректный токен авторизации", "InvalidToken", http.StatusForbidden)
	AuthorizationHeaderMissing = newCodeError(200004, "Не передан авторизационный заголовок", "AuthorizationHeaderMissing", http.StatusUnauthorized)
	InvalidAuthorizationHeader = newCodeError(200005, "Некорректный заголовок авторизации", "InvalidAuthorizationHeader", http.StatusUnauthorized)

	NotAvailableFreeSpace = newCodeError(300001, "Недостаточно места на диске", "NotAvailableFreeSpace", http.StatusBadRequest)
)

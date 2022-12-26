package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

const HTTP_HEADER_CONTENT_TYPE = "Content-Type"

const HTTP_CONTENT_TYPE_JSON_UTF8 = "application/json; charset=utf-8"
const HTTP_CONTENT_TYPE_TEXT_UTF8 = "text/plain; charset=utf-8"

func internalSetJsonResponse(aWriter http.ResponseWriter, aLogger core.Logger, aStatusCode int, aValue any, aApiError bool) {
	lValueData, lErr := json.Marshal(aValue)
	if lErr != nil {
		lMsg := fmt.Sprintf("Cannot convert %T to json", aValue)
		aLogger.Error(lMsg)

		if aApiError {
			// return empty body if we cannot convert ApiError to json
			http.Error(aWriter, "", http.StatusInternalServerError)
			return
		}

		lApiErr := ApiError{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: lMsg,
		}

		internalSetJsonResponse(aWriter, aLogger, lApiErr.ErrorCode, lApiErr, true)
		return
	}

	aWriter.Header().Set(HTTP_HEADER_CONTENT_TYPE, HTTP_CONTENT_TYPE_JSON_UTF8)
	aWriter.WriteHeader(aStatusCode)
	//TODO: handle Write() result properly
	_, _ = aWriter.Write(lValueData)
}

// func Marshal(v any) ([]byte, error) {
func SetJsonResponse(aWriter http.ResponseWriter, aValue any) {

}

func SetApiErrorResponse(aWriter http.ResponseWriter, aLogger core.Logger, aStatusCode int, aError string) {
	lApiErr := ApiError{
		ErrorCode:    aStatusCode,
		ErrorMessage: aError,
	}
	internalSetJsonResponse(aWriter, aLogger, aStatusCode, lApiErr, true)
}

func SetApiErrorResponsef(aWriter http.ResponseWriter, aLogger core.Logger, aStatusCode int, aError string, aErrorArgs ...any) {
	SetApiErrorResponse(aWriter, aLogger, aStatusCode, fmt.Sprintf(aError, aErrorArgs...))
}

func SetApiJsonResponse(aWriter http.ResponseWriter, aLogger core.Logger, aStatusCode int, aValue any) {
	internalSetJsonResponse(aWriter, aLogger, aStatusCode, aValue, false)
}

package handlers

import (
	"net/http"
	"time"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/swaggest/openapi-go/openapi3"
)

func RegisterGetArticles(aPath string, aBuilder WebServerBuilder, aFailed *bool) {
	if *aFailed {
		return
	}

	type requestParams = struct {
		FromDate time.Time `query:"from" example:"2022-11-07"`
		ToDate   time.Time `query:"to" example:"2022-11-07"`
	}

	// initialize closures variables
	lServer := aBuilder.Server()
	aBuilder.ServeMux().HandleFunc(aPath, func(w http.ResponseWriter, r *http.Request) {
		handleGetArticles(lServer, w, r)
	})

	lGet := openapi3.Operation{}
	lOpenApi := aBuilder.OpenApi()
	//TODO: handle result values properly
	_ = lOpenApi.SetRequest(&lGet, new(requestParams), http.MethodGet)
	_ = lOpenApi.SetJSONResponse(&lGet, new([]core.Article), http.StatusOK)
	_ = lOpenApi.SetJSONResponse(&lGet, ApiError{}, http.StatusBadRequest)
	_ = lOpenApi.Spec.AddOperation(http.MethodGet, aPath, lGet)
}

func handleGetArticles(aServer WebServer, aWriter http.ResponseWriter, aRequest *http.Request) {
	aServer.Log().Trace("GetArticles handler")

	lQuery := aRequest.URL.Query()

	lDefDate := time.Now().UTC()
	lDefDate = time.Date(lDefDate.Year(), lDefDate.Month(), lDefDate.Day(), 0, 0, 0, 0, lDefDate.Location())

	lFromDate, lErr := parseParamDate(lQuery.Get("from"), lDefDate.AddDate(0, 0, -7))
	if lErr != nil {
		SetApiErrorResponsef(aWriter, aServer.Log(), http.StatusBadRequest, "Invalid From parameter. %s", lErr.Error())
		return
	}
	lToDate, lErr := parseParamDate(lQuery.Get("to"), lDefDate)
	if lErr != nil {
		SetApiErrorResponsef(aWriter, aServer.Log(), http.StatusBadRequest, "Invalid To parameter. %s", lErr.Error())
		return
	}

	lArticles, lErr := aServer.DBReader().ReadArticlesByDateRange(lFromDate, lToDate)
	if lErr != nil {
		SetApiErrorResponsef(aWriter, aServer.Log(), http.StatusInternalServerError, "Cannot read articles. %s", lErr.Error())
		return
	}

	SetApiJsonResponse(aWriter, aServer.Log(), http.StatusOK, lArticles)
}

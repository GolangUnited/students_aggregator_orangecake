package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/swaggest/openapi-go/openapi3"
)

func RegisterGetArticles(aPath string, aBulder WebServerBuilder, aFailed *bool) {
	if *aFailed {
		return
	}

	type requestParams = struct {
		FromDate time.Time `query:"from" example:"2022-11-07"`
		ToDate   time.Time `query:"to" example:"2022-11-07"`
	}

	// initialize closures variables
	lServer := aBulder.Server()
	aBulder.ServeMux().HandleFunc(aPath, func(w http.ResponseWriter, r *http.Request) {
		handleGetArticles(lServer, w, r)
	})

	lGet := openapi3.Operation{}
	lOpenApi := aBulder.OpenApi()
	lOpenApi.SetRequest(&lGet, new(requestParams), http.MethodGet)
	lOpenApi.SetJSONResponse(&lGet, new([]core.Article), http.StatusOK)
	lOpenApi.SetJSONResponse(&lGet, ApiError{}, http.StatusBadRequest)
	lOpenApi.Spec.AddOperation(http.MethodGet, aPath, lGet)
}

func handleGetArticles(aServer WebServer, aWriter http.ResponseWriter, aRequest *http.Request) {
	aServer.Log().Trace("GetArticles handler")

	lQuery := aRequest.URL.Query()

	lDefDate := time.Now().UTC()
	lDefDate = time.Date(lDefDate.Year(), lDefDate.Month(), lDefDate.Day(), 0, 0, 0, 0, lDefDate.Location())

	lFromDate, lErr := parseParamDate(lQuery.Get("from"), lDefDate.AddDate(0, 0, -7))
	if lErr != nil {
		http.Error(aWriter, fmt.Sprintf("Invalid From parameter. %s", lErr.Error()), http.StatusBadRequest)
		return
	}
	lToDate, lErr := parseParamDate(lQuery.Get("to"), lDefDate)
	if lErr != nil {
		http.Error(aWriter, fmt.Sprintf("Invalid To parameter. %s", lErr.Error()), http.StatusBadRequest)
		return
	}

	//TODO: load arrticles
	lArticles := [...]core.Article{
		{Author: "Author", Title: "Title 1", Link: "LinkUrl", Description: "Summary", PublishDate: lFromDate},
		{Author: "Author", Title: "Title 2", Link: "LinkUrl", Description: "Summary", PublishDate: lFromDate},
		{Author: /*e*/ "", Title: "Title 3", Link: "LinkUrl", Description: /*em*/ "", PublishDate: lToDate},
	}

	lJson, lErr := json.Marshal(lArticles)
	if lErr != nil {
		//go:cover ignore
		http.Error(aWriter, "Cannot Marshal", http.StatusNotImplemented)
		return
	}

	aWriter.Header().Set("Content-Type", "application/json")
	aWriter.Write(lJson)
}

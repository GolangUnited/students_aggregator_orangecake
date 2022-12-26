package handlers

import (
	"fmt"
	"net/http"
)

func RegisterOpenApiJson(aPath string, aBuilder WebServerBuilder, aFailed *bool) {
	if *aFailed {
		return
	}

	// initialize closures variables
	lServer := aBuilder.Server()
	var lOpenApiJson []byte

	// Build OpenAPI document
	var lErr error
	lOpenApiJson, lErr = aBuilder.OpenApi().Spec.MarshalJSON()
	if lErr != nil {
		*aFailed = true
		lServer.Log().Error(fmt.Sprintf("Cannot generate OpenAPI document. %s", lErr.Error()))
		return
	}

	aBuilder.ServeMux().HandleFunc(aPath, func(w http.ResponseWriter, r *http.Request) {
		handleOpenApiJson(lServer, lOpenApiJson, w, r)
	})
}

func handleOpenApiJson(aServer WebServer, aOpenApiJson []byte, aWriter http.ResponseWriter, aRequest *http.Request) {
	aServer.Log().Trace("Get OpenApi Json handler")

	aWriter.Header().Set("Content-Type", "application/json")
	//TODO: handle Write() result properly
	_, _ = aWriter.Write(aOpenApiJson)
}

package handlers

import (
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

type WebServerBuilder interface {
	Server() WebServer
	OpenApi() *openapi3.Reflector
	ServeMux() *http.ServeMux
}

func NewWebServerBuilder(aServer WebServer) WebServerBuilder {
	lBuilder := builderImpl{
		server:   aServer,
		openApi:  &openapi3.Reflector{},
		serveMux: http.NewServeMux(),
	}
	lBuilder.openApi.Spec = &openapi3.Spec{Openapi: "3.0.3"}

	return lBuilder
}

type builderImpl struct {
	server   WebServer
	openApi  *openapi3.Reflector
	serveMux *http.ServeMux
}

func (b builderImpl) Server() WebServer {
	return b.server
}

func (b builderImpl) OpenApi() *openapi3.Reflector {
	return b.openApi
}

func (b builderImpl) ServeMux() *http.ServeMux {
	return b.serveMux
}

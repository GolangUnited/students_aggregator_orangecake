package handlers

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

func getErrMsg(aError error) string {
	if aError != nil {
		return aError.Error()
	}
	return ""
}

func initTestBuilder() WebServerBuilder {
	lStorage := core.MockDBStorage{}
	lWebServer, _ := NewWebServer(newMockLogger(), lStorage, lStorage)

	return NewWebServerBuilder(lWebServer)
}

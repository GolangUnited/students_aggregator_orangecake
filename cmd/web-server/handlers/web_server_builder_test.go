package handlers

import (
	"testing"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestNewWebServerBuilder(t *testing.T) {
	lSrorage := core.MockDBStorage{}
	lServer, lErr := NewWebServer(newMockLogger(), lSrorage, lSrorage)
	assert.NotNil(t, lServer, "web server")
	assert.Nil(t, lErr, "error")

	lBuilder := NewWebServerBuilder(lServer)

	assert.NotNil(t, lBuilder.Server(), "server")
	assert.NotNil(t, lBuilder.OpenApi(), "openApi")
	assert.NotNil(t, lBuilder.ServeMux(), "serveMux")
}

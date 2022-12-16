package handlers

import (
	"net/http"
)

const C_OPENAPI_UI_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta http-equiv="X-UA-Compatible" />
  <title>Demo Web API Documentation</title>
  <script src="https://unpkg.com/swagger-ui-dist@4.15.2/swagger-ui-bundle.js"></script>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@3.51.1/swagger-ui.css" />
</head>
<body>
  <div id="root"></div>
  <script>
    SwaggerUIBundle({
      dom_id: "#root",
      docExpansion: "none",
      url: "/openapi.json"});
  </script>
</body>
</html>`

func RegisterOpenApiUI(aPath string, aBuilder WebServerBuilder, aFailed *bool) {
	if *aFailed {
		return
	}

	// initialize closures variables
	lServer := aBuilder.Server()

	// Do not register Open API UI in OpenApi documentation
	aBuilder.ServeMux().HandleFunc(aPath, func(w http.ResponseWriter, r *http.Request) {
		handleOpenApiUI(lServer, w, r)
	})
}

func handleOpenApiUI(aServer WebServer, aWriter http.ResponseWriter, aRequest *http.Request) {
	aServer.Log().Trace("Get OpenApi UI handler")

	aWriter.Header().Set("Content-Type", "text/html")
	aWriter.Write([]byte(C_OPENAPI_UI_HTML))
}

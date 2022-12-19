package handlers

import (
	"strings"
	"time"
)

const (
	HTTP_DATE_PARAM_LAYOUT = "2006-01-02"
)

func parseParamDate(aParam string, aDefault time.Time) (time.Time, error) {
	aParam = strings.TrimSpace(aParam)
	if aParam == "" {
		return aDefault, nil
	}

	lDate, lErr := time.Parse(HTTP_DATE_PARAM_LAYOUT, aParam)
	if lErr != nil {
		return aDefault, lErr
	}

	return lDate, nil
}

package handlers

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseParamDate(t *testing.T) {
	type testData = struct {
		Param  string
		Result time.Time
		Error  string
	}

	lDefDate := time.Date(2001, 02, 03, 0, 0, 0, 0, time.UTC)
	lTestDate := time.Date(2010, 11, 12, 0, 0, 0, 0, time.UTC)

	lTestData := [...]testData{
		{Param: "", Result: lDefDate, Error: ""},
		{Param: "not date", Result: lDefDate, Error: "parsing time \"not date\" as \"2006-01-02\": cannot parse \"not date\" as \"2006\""},
		{Param: "2010-11-12", Result: lTestDate, Error: ""},
	}

	for i, lData := range lTestData {
		lResult, lErr := parseParamDate(lData.Param, lDefDate)
		assert.Equal(t, lData.Result, lResult, fmt.Sprintf("Date %d", i))
		assert.Equal(t, lData.Error, getErrMsg(lErr), fmt.Sprintf("Error %d", i))
	}

}

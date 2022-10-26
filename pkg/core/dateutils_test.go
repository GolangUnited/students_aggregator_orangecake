package core

import (
	"testing"
	"time"
)

func testNormalizeDate(t *testing.T, aIndex int, aTestDate, aExpectedDate string) {
	lDate, lErr := time.Parse(time.RFC3339, aTestDate)
	if lErr != nil {
		t.Error("invalid test date" + aTestDate)
	}

	lActualDate := NormalizeDate(lDate).Format(time.RFC3339)
	if aExpectedDate != lActualDate {
		t.Errorf("[%d] %s expected but %s found", aIndex+1, aExpectedDate, lActualDate)
	}
}

func TestNormalizeDate(t *testing.T) {
	type testDef struct {
		testDate string
		wantDate string
	}
	testDefs := []testDef{
		{testDate: "2022-01-02T03:04:05-06:00", wantDate: "2022-01-02T00:00:00Z"},
		{testDate: "2022-01-02T23:04:05-06:00", wantDate: "2022-01-03T00:00:00Z"},
	}

	for i, lDef := range testDefs {
		testNormalizeDate(t, i, lDef.testDate, lDef.wantDate)
	}
}

func testParseDate(t *testing.T, aIndex int, aTestDate, aExpectedDate, aError string) {

	var lError string = ""
	lDate, lErr := ParseDate("Jan _2, 2006", aTestDate)
	if lErr != nil {
		lError = lErr.Error()
	}
	if lError != aError {
		t.Errorf("[%d] '%s' error expected but '%s' found", aIndex+1, aError, lError)
	}

	lActualDate := NormalizeDate(lDate).Format(time.RFC3339)
	lExpectedDate := aExpectedDate
	if lExpectedDate == "" {
		lExpectedDate = NormalizeDate(time.Now()).Format(time.RFC3339)
	}
	if lExpectedDate != lActualDate {
		t.Errorf("[%d] %s expected but %s found", aIndex+1, lExpectedDate, lActualDate)
	}
}

func TestParseDate(t *testing.T) {
	type testDef struct {
		testDate  string
		wantDate  string
		wantError string
	}
	testDefs := []testDef{
		{testDate: "", wantDate: "", wantError: "empty Date"},
		{testDate: "02 Jan, 2022", wantDate: "", wantError: "invalid Date format Test Fail"},
		{testDate: "Jan 02, 2022", wantDate: "2022-01-02T00:00:00Z", wantError: ""},
	}

	for i, lDef := range testDefs {
		testParseDate(t, i, lDef.testDate, lDef.wantDate, lDef.wantError)
	}
}

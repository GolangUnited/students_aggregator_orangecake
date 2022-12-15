package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequiredFieldError_Error(t *testing.T) {
	testCases := map[string]struct {
		errorType error
		field     string
		result    string
	}{
		"simple input": {
			errorType: ErrFieldIsEmpty,
			field:     TitleFieldName,
			result:    "error: field is empty, field: Title",
		},
		"empty input": {
			result: "error: unknown error, field: ",
		},
	}

	for n, c := range testCases {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, c.result, RequiredFieldError{c.errorType, c.field}.Error())
		})
	}
}

func TestRequiredFieldError_Unwrap(t *testing.T) {
	testCases := map[string]struct {
		errorType error
		field     string
		result    error
	}{
		"simple input": {
			errorType: ErrFieldIsEmpty,
			field:     TitleFieldName,
			result:    ErrFieldIsEmpty,
		},
		"empty input": {
			result: ErrUnknown,
		},
	}

	for n, c := range testCases {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, c.result, RequiredFieldError{c.errorType, c.field}.Unwrap())
		})
	}
}

func TestEmptyFieldError_Error(t *testing.T) {
	testCases := map[string]struct {
		field  string
		result string
	}{
		"simple input": {
			field:  TitleFieldName,
			result: "error: field is empty, field: Title",
		},
		"empty input": {
			result: "error: field is empty, field: ",
		},
	}

	for n, c := range testCases {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, c.result, EmptyFieldError{c.field}.Error())
		})
	}
}

func TestResponseError_Error(t *testing.T) {
	testCases := map[string]struct {
		status string
		code   int
		result string
	}{
		"simple input": {
			status: "some_status",
			code:   200,
			result: "response error code: 200 status: some_status",
		},
		"empty input": {
			result: "response error code: 0 status: ",
		},
	}

	for n, c := range testCases {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, c.result, ResponseError{c.status, c.code}.Error())
		})
	}
}

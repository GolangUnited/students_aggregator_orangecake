package core

import (
	"errors"
	"fmt"
)

// Article fields name
const (
	TitleFieldName       = "Title"
	LinkFieldName        = "Link"
	AuthorFieldName      = "Author"
	DescriptionFieldName = "Description"
	PublishDateFieldName = "PublishDate"
)

// RequiredFieldError when you have any trouble with acquiring data for critical fields
// Type is field for errors ErrNodeNotFound, ErrAttributeNotExists or ErrFieldIsEmpty
type RequiredFieldError struct {
	ErrorType error
	Field     string
}

func (r RequiredFieldError) Error() string {
	if r.ErrorType == nil {
		r.ErrorType = ErrUnknown
	}
	return fmt.Sprintf("error: %s, field: %s", r.ErrorType, r.Field)
}

func (r RequiredFieldError) Unwrap() error {
	if r.ErrorType == nil {
		return ErrUnknown
	}
	return r.ErrorType
}

// Non-essential field is empty, can be sent to log
type EmptyFieldError struct {
	Field string
}

func (r EmptyFieldError) Error() string {
	return fmt.Sprintf("error: %s, field: %s", ErrFieldIsEmpty, r.Field)
}

// ResponseError when you received any error response code
type ResponseError struct {
	Status string
	Code   int
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("response error code: %d status: %s", r.Code, r.Status)
}

var (
	// ErrNodeNotFound unable to find node on page
	ErrNodeNotFound = errors.New("node not found")
	// ErrAttributeNotExists unable to find attribute in node
	ErrAttributeNotExists = errors.New("attribute doesn't exists")
	// ErrFieldIsEmpty when node is exists but there is no data
	ErrFieldIsEmpty = errors.New("field is empty")
	// ErrHTMLAccess if you can't get access to you web-page
	ErrHTMLAccess = errors.New("unable to access html page")
	// ErrNoArticles if you were unable to find any matching articles
	ErrNoArticles = errors.New("no matching articles found")
	// ErrStorageConnection cannot connect to storage
	ErrStorageConnection = errors.New("unable to connect to storage")
	// ErrEmptyDate empty date
	ErrEmptyDate = errors.New("empty Date")
	// ErrInvalidDateFormat invalid data format
	ErrInvalidDateFormat = errors.New("invalid Date format")
	//ErrUrlVisit scrapper cant visit url
	ErrUrlVisit = errors.New("unable to visit URL")
	// ErrConfigFileIsEmpty empty config file
	ErrConfigFileIsEmpty = errors.New("config file is empty")
	// ErrEmptyConfig config struct is empty
	ErrEmptyConfig = errors.New("empty config struct")
	// ErrEmptyEnvVariable env variable not found
	ErrEmptyEnvVariable = errors.New("env variable not found")
	// ErrNoHandlersInConfig unable to unmarshal info about handlers from yaml file
	ErrNoHandlersInConfig = errors.New("unable to get handlers info from config")
	// ErrUnknown unknown error
	ErrUnknown = errors.New("unknown error")
)

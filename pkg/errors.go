package pkg

import (
	"errors"
)

var (
	// ErrRequiredFieldNotFound use this error when you can't get data for Title or URL
	ErrRequiredFieldNotFound = errors.New("data for required field not found")

	// ErrFieldNotFound use this error when you can't get data for Author, Date or Description
	ErrFieldNotFound = errors.New("data for field not found")

	// ErrHTMLAccess return this error from GetArticles if you can't get access to you web-page otherwise return nil error and articles which you were able to get
	ErrHTMLAccess = errors.New("unable to access page")

	ErrEmptyDate = errors.New("empty Date")

	ErrInvalidDateFormat = errors.New("invalid Date format")
)

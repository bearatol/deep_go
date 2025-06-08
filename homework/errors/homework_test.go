package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	str := fmt.Sprintf("%d errors occured:\n", len(e.Errors))
	for _, v := range e.Errors {
		str += "\t* " + v.Error()
	}
	return str + "\n"
}

func Append(err error, errs ...error) *MultiError {
	mErr, ok := err.(*MultiError)
	if !ok {
		mErr = &MultiError{}
	}
	mErr.Errors = append(mErr.Errors, errs...)
	return mErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

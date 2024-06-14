package errors

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestErrorL10nTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorL10nTestSuite))
}

type ErrorL10nTestSuite struct {
	suite.Suite
}

func (s *ErrorL10nTestSuite) SetupSuite() {
}

func (s *ErrorL10nTestSuite) Test_ConstructNewError() {
	err := NewI18nError("some-key")
	s.Assert().Equal("some-key", err.Error())
}

func (s *ErrorL10nTestSuite) Test_MustImplementOfStdError() {
	var v interface{} = NewI18nError("some-key")
	_, ok := v.(error)
	s.Assert().True(ok, "New error must implement of standard error")
}

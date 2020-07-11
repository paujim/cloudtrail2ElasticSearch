package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockESClient struct {
	mock.Mock
}

func (m *MockESClient) Do(input *http.Request) (*http.Response, error) {
	args := m.Called(input)
	var resp *http.Response
	if args.Get(0) != nil {
		resp = args.Get(0).(*http.Response)
	}
	return resp, args.Error(1)
}

func TestPost(t *testing.T) {

	t.Run("Post data return 201", func(t *testing.T) {

		mockClient := &MockESClient{}
		mockClient.On("Do", mock.Anything).Return(&http.Response{StatusCode: 201}, nil)

		lambdaHandler := &LambdaHandler{mockClient}

		err := lambdaHandler.Handle(nil, CloudTrailEvent{})
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})
	t.Run("Post data return error", func(t *testing.T) {

		mockClient := &MockESClient{}
		mockClient.On("Do", mock.Anything).Return(nil, errors.New("BAD REQUEST"))

		lambdaHandler := &LambdaHandler{mockClient}

		err := lambdaHandler.Handle(nil, CloudTrailEvent{})
		assert.Error(t, err, "BAD REQUEST")
		mockClient.AssertExpectations(t)
	})
}

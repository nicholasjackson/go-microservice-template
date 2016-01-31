package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockStatsD struct {
	mock.Mock
}

func (m *MockStatsD) Increment(label string) {
	_ = m.Mock.Called(label)
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = m.Mock.Called(w, r)
}

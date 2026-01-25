package errmap

import (
	"errors"

	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
)

type ErrorRegistry interface {
	Lookup(err error) (Result, bool)
}
type Result struct {
	Status   int
	Response errorResponses.ErrorResponse
}

type Registry struct {
	mappings map[error]Result
}

func NewRegistry() *Registry {
	return &Registry{
		mappings: make(map[error]Result),
	}
}

func (m *Registry) Register(err error, status int, response errorResponses.ErrorResponse) {
	m.mappings[err] = Result{
		Status:   status,
		Response: response,
	}
}

func (m *Registry) Lookup(err error) (Result, bool) {
	// Traversing is intentional not accidental. Wrapped errors may come.
	for e, mapping := range m.mappings {
		if errors.Is(err, e) {
			return mapping, true
		}
	}
	return Result{}, false
}

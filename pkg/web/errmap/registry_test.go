package errmap

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/greenbone/opensight-golang-libraries/pkg/errorResponses"
	"github.com/stretchr/testify/assert"
)

func TestRegistry_Lookup(t *testing.T) {
	reg := NewRegistry()

	errSentinel := errors.New("sentinel error")
	errNotRegistered := errors.New("not registered")

	mockResponse := errorResponses.ErrorInternalResponse
	reg.Register(errSentinel, http.StatusBadRequest, mockResponse)

	tests := []struct {
		name       string
		inputErr   error
		wantFound  bool
		wantStatus int
	}{
		{
			name:       "Direct match returns correctly",
			inputErr:   errSentinel,
			wantFound:  true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Wrapped error match (using %w) returns correctly",
			inputErr:   fmt.Errorf("context: %w", errSentinel),
			wantFound:  true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:      "Unregistered error returns false",
			inputErr:  errNotRegistered,
			wantFound: false,
		},
		{
			name:      "Nil error returns false",
			inputErr:  nil,
			wantFound: false,
		},
	}

	// 3. Execution
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := reg.Lookup(tt.inputErr)
			assert.Equal(t, tt.wantFound, found)
			if tt.wantFound {
				assert.Equal(t, tt.wantStatus, result.Status)
			}
		})
	}
}

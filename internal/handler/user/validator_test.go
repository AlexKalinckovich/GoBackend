package user

import (
	"testing"

	"github.com/brota/gobackend/internal/custom_errors/validation"
	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

func TestValidator_ValidateCreateOrUpdate(t *testing.T) {
	v := NewUserValidator()

	tests := []struct {
		name      string
		req       CreateUserRequest
		wantError bool
		errCount  int
	}{
		{
			name:      "Valid",
			req:       CreateUserRequest{Name: "John", Surname: "Doe", Age: 25},
			wantError: false,
		},
		{
			name:      "Empty Name",
			req:       CreateUserRequest{Name: "   ", Surname: "Doe", Age: 25},
			wantError: true,
			errCount:  1,
		},
		{
			name:      "Invalid Age",
			req:       CreateUserRequest{Name: "John", Surname: "Doe", Age: -5},
			wantError: true,
			errCount:  1,
		},
		{
			name:      "Multiple Errors",
			req:       CreateUserRequest{Name: "", Surname: "", Age: 200},
			wantError: true,
			errCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errCh := make(chan *validation.AggregateError)
			go v.ValidateCreateOrUpdate(tt.req, errCh)
			err := <-errCh

			if tt.wantError {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errCount, len(err.Errors))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidator_ValidatePatch(t *testing.T) {
	v := NewUserValidator()

	tests := []struct {
		name      string
		req       PatchUserRequest
		wantError bool
		errCount  int
	}{
		{
			name:      "Valid Full Patch",
			req:       PatchUserRequest{Name: ptr("John"), Surname: ptr("Doe"), Age: ptr(25)},
			wantError: false,
		},
		{
			name:      "Valid Partial Patch",
			req:       PatchUserRequest{Age: ptr(30)},
			wantError: false,
		},
		{
			name:      "Invalid Name",
			req:       PatchUserRequest{Name: ptr("")},
			wantError: true,
			errCount:  1,
		},
		{
			name:      "Multiple Errors",
			req:       PatchUserRequest{Name: ptr(" "), Age: ptr(-1)},
			wantError: true,
			errCount:  2,
		},
		{
			name:      "Empty Request",
			req:       PatchUserRequest{},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errCh := make(chan *validation.AggregateError)
			go v.ValidatePatch(tt.req, errCh)
			err := <-errCh

			if tt.wantError {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errCount, len(err.Errors))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

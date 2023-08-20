package services

import (
	"marathon-postgresql/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRunner(t *testing.T) {
	tests := []struct {
		name   string
		runner *models.Runner
		want   *models.ResponseError
	}{
		{
			name: "Invalid_first_name",
			runner: &models.Runner{
				LastName: "Smith",
				Age:      30,
				Country:  "United States",
			},
			want: &models.ResponseError{
				Message: "Invalid first name",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_last_name",
			runner: &models.Runner{
				FirstName: "John",
				Age:       30,
				Country:   "United States",
			},
			want: &models.ResponseError{
				Message: "Invalid last name",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_age",
			runner: &models.Runner{
				FirstName: "John",
				LastName:  "Smith",
				Age:       3000,
				Country:   "United States",
			},
			want: &models.ResponseError{
				Message: "Invalid age",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_country",
			runner: &models.Runner{
				FirstName: "John",
				LastName:  "Smith",
				Age:       30,
			},
			want: &models.ResponseError{
				Message: "Invalid country",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Valid_runner",
			runner: &models.Runner{
				FirstName: "John",
				LastName:  "Smith",
				Age:       30,
				Country:   "United States",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			responseErr := validateRunner(test.runner)
			assert.Equal(t, test.want, responseErr)
		})
	}
}

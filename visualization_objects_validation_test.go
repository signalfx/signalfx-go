package signalfx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisualizationObjectsValidationValidate(t *testing.T) {
	tests := []struct {
		input         VisualizationObjectsValidation
		expectedError bool
	}{
		{FULL, false},
		{TERRAFORM, false},
		{"INVALID", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectedError {
				assert.Error(t, err, "Expected error for invalid input.")
			} else {
				assert.NoError(t, err, "Unexpected error for valid input.")
			}
		})
	}
}

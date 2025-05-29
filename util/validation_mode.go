package util

import "errors"

type ValidationMode string

const (
	FULL      ValidationMode = "FULL"
	TERRAFORM ValidationMode = "TERRAFORM"
)

func ValidateValidationMode(validationMode ValidationMode) error {
	switch validationMode {
	case FULL, TERRAFORM:
		return nil
	default:
		return errors.New("invalid validation mode: " + string(validationMode))
	}
}

package signalfx

import "errors"

type VisualizationObjectsValidation string

const (
	FULL      VisualizationObjectsValidation = "FULL"
	TERRAFORM VisualizationObjectsValidation = "TERRAFORM"
)

func (v VisualizationObjectsValidation) Validate() error {
	switch v {
	case FULL, TERRAFORM:
		return nil
	default:
		return errors.New("invalid validation mode: " + string(v))
	}
}

package slo

import (
	"encoding/json"
	"fmt"
)

const (
	RollingWindowTarget  = "RollingWindow"
	CalendarWindowTarget = "CalendarWindow"
)

type SloTarget struct {
	BaseSloTarget
	*RollingWindowSloTarget
	*CalendarWindowSloTarget
}

type BaseSloTarget struct {
	Slo           float64        `json:"slo,omitempty"`
	SloAlertRules []SloAlertRule `json:"sloAlertRules,omitempty"`
	Type          string         `json:"type,omitempty"`
}

type RollingWindowSloTarget struct {
	CompliancePeriod string `json:"compliancePeriod,omitempty"`
}

type CalendarWindowSloTarget struct {
	CycleType  string `json:"cycleType,omitempty"`
	CycleStart string `json:"cycleStart,omitempty"`
}

func (target *SloTarget) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &target.BaseSloTarget); err != nil {
		return err
	}
	switch target.Type {
	case RollingWindowTarget:
		target.RollingWindowSloTarget = &RollingWindowSloTarget{}
		return json.Unmarshal(data, target.RollingWindowSloTarget)
	case CalendarWindowTarget:
		target.CalendarWindowSloTarget = &CalendarWindowSloTarget{}
		return json.Unmarshal(data, target.CalendarWindowSloTarget)
	default:
		return fmt.Errorf("unrecognized SLO target type %s", target.Type)
	}
}

func (target *SloTarget) MarshalJSON() ([]byte, error) {
	switch target.Type {
	case RollingWindowTarget:
		return json.Marshal(struct {
			BaseSloTarget
			*RollingWindowSloTarget
		}{target.BaseSloTarget, target.RollingWindowSloTarget})
	case CalendarWindowTarget:
		return json.Marshal(struct {
			BaseSloTarget
			*CalendarWindowSloTarget
		}{target.BaseSloTarget, target.CalendarWindowSloTarget})
	default:
		return nil, fmt.Errorf("unrecognized SLO target type %s", target.Type)
	}
}

package slo

import (
	"encoding/json"
	"fmt"
	"github.com/signalfx/signalfx-go/detector"
)

const (
	BreachRule          = "BREACH"
	ErrorBudgetLeftRule = "ERROR_BUDGET_LEFT"
	BurnRateRule        = "BURN_RATE"
)

type BreachSloAlertRule struct {
	Rules []*BreachDetectorRule `json:"rules,omitempty"`
}

type BreachDetectorRule struct {
	detector.Rule
	Parameters *BreachDetectorParameters `json:"parameters,omitempty"`
}

type BreachDetectorParameters struct {
	FireLasting      string  `json:"fireLasting,omitempty"`
	PercentOfLasting float64 `json:"percentOfLasting,omitempty"`
}

type ErrorBudgetLeftSloAlertRule struct {
	Rules []*ErrorBudgetLeftDetectorRule `json:"rules,omitempty"`
}

type ErrorBudgetLeftDetectorRule struct {
	detector.Rule
	Parameters *ErrorBudgetLeftDetectorParameters `json:"parameters,omitempty"`
}

type ErrorBudgetLeftDetectorParameters struct {
	FireLasting            string  `json:"fireLasting,omitempty"`
	PercentOfLasting       float64 `json:"percentOfLasting,omitempty"`
	PercentErrorBudgetLeft float64 `json:"percentErrorBudgetLeft,omitempty"`
}

type BurnRateSloAlertRule struct {
	Rules []*BurnRateDetectorRule `json:"rules,omitempty"`
}

type BurnRateDetectorRule struct {
	detector.Rule
	Parameters *BurnRateDetectorParameters `json:"parameters,omitempty"`
}

type BurnRateDetectorParameters struct {
	ShortWindow1       string  `json:"shortWindow1,omitempty"`
	LongWindow1        string  `json:"longWindow1,omitempty"`
	ShortWindow2       string  `json:"shortWindow2,omitempty"`
	LongWindow2        string  `json:"longWindow2,omitempty"`
	BurnRateThreshold1 float64 `json:"burnRateThreshold1,omitempty"`
	BurnRateThreshold2 float64 `json:"burnRateThreshold2,omitempty"`
}

type BaseSloAlertRule struct {
	Type string `json:"type,omitempty"`
}

type SloAlertRule struct {
	BaseSloAlertRule
	*BreachSloAlertRule
	*ErrorBudgetLeftSloAlertRule
	*BurnRateSloAlertRule
}

func (rule *SloAlertRule) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &rule.BaseSloAlertRule); err != nil {
		return err
	}
	switch rule.Type {
	case BreachRule:
		rule.BreachSloAlertRule = &BreachSloAlertRule{}
		return json.Unmarshal(data, rule.BreachSloAlertRule)
	case ErrorBudgetLeftRule:
		rule.ErrorBudgetLeftSloAlertRule = &ErrorBudgetLeftSloAlertRule{}
		return json.Unmarshal(data, rule.ErrorBudgetLeftSloAlertRule)
	case BurnRateRule:
		rule.BurnRateSloAlertRule = &BurnRateSloAlertRule{}
		return json.Unmarshal(data, rule.BurnRateSloAlertRule)
	default:
		return fmt.Errorf("unrecognized SLO alert rule type %s", rule.Type)
	}
}

func (rule *SloAlertRule) MarshalJSON() ([]byte, error) {
	switch rule.Type {
	case BreachRule:
		return json.Marshal(struct {
			BaseSloAlertRule
			*BreachSloAlertRule
		}{rule.BaseSloAlertRule, rule.BreachSloAlertRule})
	case ErrorBudgetLeftRule:
		return json.Marshal(struct {
			BaseSloAlertRule
			*ErrorBudgetLeftSloAlertRule
		}{rule.BaseSloAlertRule, rule.ErrorBudgetLeftSloAlertRule})
	case BurnRateRule:
		return json.Marshal(struct {
			BaseSloAlertRule
			*BurnRateSloAlertRule
		}{rule.BaseSloAlertRule, rule.BurnRateSloAlertRule})
	default:
		return nil, fmt.Errorf("unrecognized SLO alert rule type %s", rule.Type)
	}
}

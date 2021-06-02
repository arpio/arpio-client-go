package arpio

import (
	"encoding/json"
	"fmt"
	"time"
)

// Enumeration of the types of Arpio applications.
const (
	StandardAppType  = "standard"
	TerraformAppType = "terraform"
)

type App struct {
	AccountID          string          `json:"accountId"`
	AppID              string          `json:"appId,omitempty"`
	AppType            string          `json:"type"`
	CreatedAt          time.Time       `json:"createdAt,omitempty"`
	Name               string          `json:"name"`
	NotificationEmails []string        `json:"notificationEmails,omitempty"`
	RPO                int             `json:"rpo"`
	SelectionRules     []SelectionRule `json:"-"`
	SourceAwsAccountID string          `json:"sourceAwsAccountId"`
	SourceRegion       string          `json:"sourceRegion"`
	SyncPhase          string          `json:"syncPhase,omitempty"`
	TargetAwsAccountID string          `json:"targetAwsAccountId"`
	TargetRegion       string          `json:"targetRegion"`

	// JSON serialization helpers
	RawSelectionRules []json.RawMessage `json:"selectionRules"`
}

func (a App) MarshalJSON() ([]byte, error) {
	// Use a type alias to avoid invoking this function recursively
	type AppDTO App

	// Marshal SelectionRules through RawSelectionRules
	a.RawSelectionRules = []json.RawMessage{}
	if a.SelectionRules != nil {
		for _, r := range a.SelectionRules {
			b, err := json.Marshal(r)
			if err != nil {
				return nil, err
			}
			a.RawSelectionRules = append(a.RawSelectionRules, b)
		}
	}

	return json.Marshal((AppDTO)(a))
}

func (a *App) UnmarshalJSON(b []byte) error {
	// Use a type alias to avoid invoking this function recursively
	type AppDTO App
	err := json.Unmarshal(b, (*AppDTO)(a))
	if err != nil {
		return err
	}

	// Unmarshal RawSelectionRules into SelectionRules
	a.SelectionRules = []SelectionRule{}
	for _, rawSelectionRule := range a.RawSelectionRules {
		var baseRule selectionRule
		err = json.Unmarshal(rawSelectionRule, &baseRule)
		if err != nil {
			return err
		}

		var rule SelectionRule
		switch baseRule.RuleType {
		case ArnRuleType:
			r := &ArnRule{}
			err = json.Unmarshal(rawSelectionRule, r)
			rule = *r
		case TagRuleType:
			r := &TagRule{}
			err = json.Unmarshal(rawSelectionRule, r)
			rule = *r
		default:
			panic(fmt.Sprintf("unhandled selection rule type: %s", baseRule.RuleType))
		}
		if err != nil {
			return err
		}
		a.SelectionRules = append(a.SelectionRules, rule)
	}

	return nil
}

// SyncPair returns a SyncPair struct for the App's endpoint information.
func (a App) SyncPair() SyncPair {
	return NewSyncPair(
		a.SourceAwsAccountID,
		a.SourceRegion,
		a.TargetAwsAccountID,
		a.TargetRegion,
	)
}

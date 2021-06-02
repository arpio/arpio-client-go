package arpio

const (
	ArnRuleType = "arn"
	TagRuleType = "tag"
)

type SelectionRule interface {
	GetRuleType() string
}

// Embeddable struct to assist with polymorphic unmarshaling of rule types.
type selectionRule struct {
	RuleType string `json:"ruleType"`
}

// GetRuleType returns the rule type.
func (r selectionRule) GetRuleType() string {
	return r.RuleType
}

type ArnRule struct {
	selectionRule
	Arns []string `json:"arns"`
}

// NewArnRule creates an ArnRule that will match resources with the
// specified ARN strings.
func NewArnRule(arns []string) ArnRule {
	return ArnRule{
		selectionRule: selectionRule{RuleType: ArnRuleType},
		Arns:          arns,
	}
}

type TagRule struct {
	selectionRule
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewTagRule creates an TagRule that will match resources with the
// specified tag name and value (which may be empty).
func NewTagRule(name, value string) TagRule {
	return TagRule{
		selectionRule: selectionRule{RuleType: TagRuleType},
		Name:          name,
		Value:         value,
	}
}

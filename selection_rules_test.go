package arpio

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestArnRuleJSON(t *testing.T) {
	orig := ArnRule{
		selectionRule: selectionRule{
			RuleType: ArnRuleType,
		},
		Arns: []string{"arn:a", "arn:b"},
	}
	bytes, err := json.Marshal(orig)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"ruleType":"arn","arns":["arn:a","arn:b"]}`
	if string(bytes) != expected {
		t.Fatalf("%s != %s", bytes, expected)
	}

	var decoded ArnRule
	err = json.Unmarshal(bytes, &decoded)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Fatalf("%v != %v", orig, decoded)
	}
}

func TestTagRuleJSON(t *testing.T) {
	orig := TagRule{
		selectionRule: selectionRule{
			RuleType: TagRuleType,
		},
		Name:  "foo",
		Value: "bar",
	}
	bytes, err := json.Marshal(orig)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"ruleType":"tag","name":"foo","value":"bar"}`
	if string(bytes) != expected {
		t.Fatalf("%s != %s", bytes, expected)
	}

	var decoded TagRule
	err = json.Unmarshal(bytes, &decoded)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Fatalf("%v != %v", orig, decoded)
	}
}

package arpio

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestAppJSON(t *testing.T) {
	orig := App{
		AccountID:          "a",
		AppID:              "b",
		AppType:            "c",
		CreatedAt:          time.Time{},
		Name:               "d",
		NotificationEmails: []string{"a@example.com", "b@example.com"},
		RPO:                0,
		SelectionRules: []SelectionRule{
			NewArnRule([]string{"arn:a", "arn:b"}),
			NewTagRule("foo", "bar"),
			NewTagRule("foo2", "bar2"),
		},
		SourceAwsAccountID: "e",
		SourceRegion:       "f",
		SyncPhase:          "g",
		TargetAwsAccountID: "h",
		TargetRegion:       "i",
	}
	bytes, err := json.Marshal(orig)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"accountId":"a","appId":"b","type":"c","createdAt":"0001-01-01T00:00:00Z","name":"d","notificationEmails":["a@example.com","b@example.com"],"rpo":0,"sourceAwsAccountId":"e","sourceRegion":"f","syncPhase":"g","targetAwsAccountId":"h","targetRegion":"i","selectionRules":[{"ruleType":"arn","arns":["arn:a","arn:b"]},{"ruleType":"tag","name":"foo","value":"bar"},{"ruleType":"tag","name":"foo2","value":"bar2"}]}`
	if string(bytes) != expected {
		t.Fatalf("%s != %s", bytes, expected)
	}

	var decoded App
	err = json.Unmarshal(bytes, &decoded)
	if err != nil {
		t.Fatal(err)
	}

	// Clear raw fields that may be dirty from serialization
	decoded.RawSelectionRules = nil
	orig.RawSelectionRules = nil

	if !reflect.DeepEqual(orig, decoded) {
		t.Fatalf("%v != %v", orig, decoded)
	}
}

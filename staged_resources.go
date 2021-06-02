package arpio

import (
	"encoding/json"
	"log"
)

const (
	BackupRecoveryPointExtraType  = "backupRecoveryPoint"
	BackupVaultExtraType          = "backupVault"
	EC2ImageExtraType             = "ec2Image"
	EC2SnapshotExtraType          = "ec2Snapshot"
	KMSKeyExtraType               = "kmsKey"
	RDSDBClusterSnapshotExtraType = "rdsDbClusterSnapshot"
	RDSDBSnapshotExtraType        = "rdsDbSnapshot"
	RDSOptionGroupExtraType       = "rdsOptionGroup"
)

const (
	SourceEnvironment = "source"
	TargetEnvironment = "target"
)

type StagedExtra interface {
	GetType() string
	GetEnvironment() string
}

// Embeddable struct to assist with polymorphic unmarshaling of staged extra types.
type stagedExtra struct {
	Type        string `json:"type"`
	Environment string `json:"environment"`
}

// GetType returns the staged extra type.
func (e stagedExtra) GetType() string {
	return e.Type
}

// GetEnvironment returns the staged extra environment.
func (e stagedExtra) GetEnvironment() string {
	return e.Environment
}

type BackupVaultExtra struct {
	stagedExtra
	BackupVaultName string `json:"backupVaultName"`
}

type BackupRecoveryPointExtra struct {
	stagedExtra
	BackupVaultName  string            `json:"backupVaultName"`
	RecoveryPointARN string            `json:"recoveryPointArn"`
	RestoreMetadata  map[string]string `json:"restoreMetadata"`
}

type EC2ImageExtra struct {
	stagedExtra
	ImageARN     string   `json:"imageArn"`
	SnapshotARNs []string `json:"snapshotArns"`
}

type EC2SnapshotExtra struct {
	stagedExtra
	SnapshotARN string `json:"snapshotArn"`
}

type KMSKeyExtra struct {
	stagedExtra
	KMSKeyARN string `json:"kmsKeyArn"`
}

type RDSDBClusterSnapshotExtra struct {
	stagedExtra
	DBClusterSnapshotARN string `json:"dbClusterSnapshotArn"`
}

type RDSDBSnapshotExtra struct {
	stagedExtra
	DBSnapshotARN string `json:"dbSnapshotArn"`
}

type RDSOptionGroupExtra struct {
	stagedExtra
	OptionGroupARN string `json:"optionGroupArn"`
}

type StagedResource struct {
	ARN    string            `json:"arn"`
	Tags   map[string]string `json:"tags"`
	Type   string            `json:"type"`
	Extras []StagedExtra     `json:"-"`

	// JSON serialization helpers
	RawExtras []json.RawMessage `json:"extras"`
}

func (sr StagedResource) MarshalJSON() ([]byte, error) {
	// Use a type alias to avoid invoking this function recursively
	type StagedResourceDTO StagedResource

	// Marshal Extras through RawExtras
	sr.RawExtras = []json.RawMessage{}
	if sr.Extras != nil {
		for _, e := range sr.Extras {
			b, err := json.Marshal(e)
			if err != nil {
				return nil, err
			}
			sr.RawExtras = append(sr.RawExtras, b)
		}
	}

	return json.Marshal((StagedResourceDTO)(sr))
}

func (sr *StagedResource) UnmarshalJSON(b []byte) error {
	// Use a type alias to avoid invoking this function recursively
	type StagedResourceDTO StagedResource
	err := json.Unmarshal(b, (*StagedResourceDTO)(sr))
	if err != nil {
		return err
	}

	// Unmarshal RawExtras into Extras
	sr.Extras = []StagedExtra{}
	for _, rawExtra := range sr.RawExtras {
		var baseExtra stagedExtra
		err = json.Unmarshal(rawExtra, &baseExtra)
		if err != nil {
			return err
		}

		var extra StagedExtra
		switch baseExtra.Type {
		case BackupRecoveryPointExtraType:
			e := &BackupRecoveryPointExtra{}
			err = json.Unmarshal(rawExtra, e)
			extra = *e
		case BackupVaultExtraType:
			e := &BackupVaultExtra{}
			err = json.Unmarshal(rawExtra, e)
			extra = *e
		case EC2ImageExtraType:
			e := &EC2ImageExtra{}
			err = json.Unmarshal(rawExtra, e)
			extra = *e
		case EC2SnapshotExtraType:
			e := &EC2SnapshotExtra{}
			err = json.Unmarshal(rawExtra, e)
			extra = *e
		case KMSKeyExtraType:
			e := &KMSKeyExtra{}
			err = json.Unmarshal(rawExtra, e)
			extra = *e
		case RDSDBClusterSnapshotExtraType:
			e := &RDSDBClusterSnapshotExtra{}
			err = json.Unmarshal(rawExtra, e)
			extra = *e
		case RDSDBSnapshotExtraType:
			e := &RDSDBSnapshotExtra{}
			err = json.Unmarshal(rawExtra, e)
			extra = *e
		case RDSOptionGroupExtraType:
			e := &RDSOptionGroupExtra{}
			err = json.Unmarshal(rawExtra, e)
			extra = *e
		default:
			log.Printf("[DEBUG] Ignoring extra with unknown type: %s", baseExtra.Type)
		}
		if err != nil {
			return err
		}
		sr.Extras = append(sr.Extras, extra)
	}

	return nil
}

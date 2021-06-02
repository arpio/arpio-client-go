package arpio

import (
	"fmt"
)

type SyncPair struct {
	Source SyncEndpoint
	Target SyncEndpoint
}

// NewSyncPair creates a SyncPair with the specified endpoint information.
func NewSyncPair(sourceAWSAccountID, sourceRegion, targetAWSAccountID, targetRegion string) SyncPair {
	return SyncPair{
		Source: SyncEndpoint{sourceAWSAccountID, sourceRegion},
		Target: SyncEndpoint{targetAWSAccountID, targetRegion},
	}
}

func (sp SyncPair) String() string {
	return fmt.Sprintf("%s/%s", sp.Source, sp.Target)
}

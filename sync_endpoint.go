package arpio

import (
	"fmt"
)

type SyncEndpoint struct {
	AccountID string
	Region    string
}

func (ep SyncEndpoint) String() string {
	return fmt.Sprintf("%s/%s", ep.AccountID, ep.Region)
}

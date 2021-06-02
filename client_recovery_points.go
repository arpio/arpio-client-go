package arpio

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const RecoveryPointPollPeriod = 5 * time.Second

// ListRecoveryPoints lists all the recovery points for the specified sync pair.
// If either timestampStart or timestampEnd is nil, that timestamp is
// unconstrained in that direction.
func (c *Client) ListRecoveryPoints(syncPair SyncPair, timestampStart, timestampEnd *time.Time) (rps []RecoveryPoint, err error) {
	spURL := c.syncPairPath(syncPair)

	v := url.Values{}
	if timestampStart != nil {
		v["timestampStart"] = []string{timestampStart.Format(time.RFC3339)}
	}
	if timestampEnd != nil {
		v["timestampEnd"] = []string{timestampEnd.Format(time.RFC3339)}
	}
	query := ""
	if len(v) > 0 {
		query = "?" + v.Encode()
	}

	u := fmt.Sprintf("%s/recoveryPoints%s", spURL, query)

	_, err = c.apiGet(u, &rps)
	if err != nil {
		return rps, err
	}

	return rps, nil
}

// GetRecoveryPoint gets the recovery point with the specified ID.
func (c *Client) GetRecoveryPoint(syncPair SyncPair, recoveryPointID string) (rp *RecoveryPoint, err error) {
	u := c.recoveryPointPath(syncPair, recoveryPointID)

	var r RecoveryPoint
	status, err := c.apiGet(u, &r)
	if status == http.StatusNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// UpdateRecoveryPoint updates the mutable properties of the specified recovery point.
func (c *Client) UpdateRecoveryPoint(syncPair SyncPair, rp RecoveryPoint) (updated RecoveryPoint, err error) {
	u := c.recoveryPointPath(syncPair, rp.RecoveryPointID)

	_, err = c.apiPut(u, rp, &updated)
	if err != nil {
		return updated, err
	}

	return updated, nil
}

// ListRecoveryPointResources lists all the staged resources in the specified recovery point.
func (c *Client) ListRecoveryPointResources(syncPair SyncPair, rp RecoveryPoint) (srs []StagedResource, err error) {
	rpURL := c.recoveryPointPath(syncPair, rp.RecoveryPointID)
	u := fmt.Sprintf("%s/resources", rpURL)

	_, err = c.apiGet(u, &srs)
	if err != nil {
		return srs, err
	}

	return srs, nil
}

// FindLatestRecoveryPoint finds the most recent recovery point that
// matches the timestamp criteria.  If timestampMin is empty, the minimum
// is not constrained.  If timestampMax is empty, the maximum time is not
// constrained. If no recovery points exist between the time constraints,
// nil is returned.
func (c *Client) FindLatestRecoveryPoint(syncPair SyncPair, timestampMin, timestampMax *time.Time) (rp *RecoveryPoint, err error) {
	// List all recovery points
	allRPs, err := c.ListRecoveryPoints(syncPair, timestampMin, timestampMax)
	if err != nil {
		return rp, err
	}

	// Select the most recent recovery point
	for i, r := range allRPs {
		if rp == nil || r.Timestamp.After(rp.Timestamp) {
			rp = &allRPs[i]
		}
	}

	return rp, nil
}

// MustFindLatestRecoveryPoint finds the most recent recovery point that
// matches the timestamp criteria. If timeout is > 0, the function tries to
// find a matching recovery point until the timeout has elapsed.  An error is
// returned if no matching recovery point could be found.
func (c *Client) MustFindLatestRecoveryPoint(syncPair SyncPair, timestampMin, timestampMax *time.Time, timeout time.Duration) (rp *RecoveryPoint, err error) {
	const zeroDuration = time.Duration(0)
	for timeoutAt := time.Now().Add(timeout); timeout == zeroDuration || time.Now().Before(timeoutAt); {
		rp, err = c.FindLatestRecoveryPoint(syncPair, timestampMin, timestampMax)
		if err != nil {
			return rp, err
		}
		if rp != nil || timeout == zeroDuration {
			break
		}
		log.Printf("[DEBUG] Waiting for a matching recovery point to exist")
		time.Sleep(RecoveryPointPollPeriod)
	}

	// If we didn't find a recovery point, prepare an error
	if rp == nil {
		if timestampMin != nil && timestampMax != nil {
			err = fmt.Errorf("there are no recovery points between "+
				"%q and %q; change timestamp_min to an earlier time or "+
				"remove it from your config to use an older recovery point, "+
				"or remove timestamp from your config to use the most recent "+
				"recovery point",
				timestampMin.Format(time.RFC3339),
				timestampMax.Format(time.RFC3339))
		} else if timestampMin != nil {
			err = fmt.Errorf("there are no recovery points on or "+
				"after %q; change timestamp_min to an earlier time or remove "+
				"it from your config to use an older recovery point",
				timestampMin.Format(time.RFC3339))
		} else if timestampMax != nil {
			err = fmt.Errorf("there are no recovery points on or "+
				"before %q; change timestamp to a later time or remove "+
				"it from your config to use a newer recovery point",
				timestampMax.Format(time.RFC3339))
		} else {
			// No timestampMax specified, but no recovery points at all
			err = fmt.Errorf("no recovery points exist for this " +
				"application yet; please wait for the first recovery point " +
				"to be created")
		}
	}

	return rp, err
}

// ProtectRecoveryPoint sets the "Protected" attribute to true and updates the
// recovery point in the Arpio service.
func (c *Client) ProtectRecoveryPoint(syncPair SyncPair, recoveryPoint RecoveryPoint) (protected RecoveryPoint, err error) {
	if recoveryPoint.Protected {
		return recoveryPoint, nil
	}

	protected = recoveryPoint
	protected.Protected = true

	protected, err = c.UpdateRecoveryPoint(syncPair, protected)
	if err != nil {
		return protected, err
	}

	return protected, nil
}

func (c *Client) syncPairPath(syncPair SyncPair) string {
	return fmt.Sprintf(
		"/accounts/%s/syncPairs/%s/%s/%s/%s",
		c.AccountID,
		syncPair.Source.AccountID,
		syncPair.Source.Region,
		syncPair.Target.AccountID,
		syncPair.Target.Region,
	)
}

func (c *Client) recoveryPointPath(syncPair SyncPair, recoveryPointID string) string {
	spURL := c.syncPairPath(syncPair)
	return fmt.Sprintf(
		"%s/recoveryPoints/%s",
		spURL,
		recoveryPointID,
	)
}

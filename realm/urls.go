// package realm contains helpers for constructing realm-specific urls and
// config.
package realm

import (
	"fmt"
	"strings"
)

// IngestURLForRealm returns the base ingest URL for a particular SignalFx
// realm
func IngestURLForRealm(realm string) string {
	return fmt.Sprintf("https://ingest.%s.signalfx.com", realm)
}

// APIURLForRealm returns the base API URL for a particular SignalFx realm
func APIURLForRealm(realm string) string {
	return fmt.Sprintf("https://api.%s.signalfx.com", realm)
}

// DatapointEndpointForRealm returns the endpoint to which datapoints should be
// POSTed for a particular realm.
func DatapointEndpointForRealm(realm string) string {
	return DatapointEndpointForIngestURL(IngestURLForRealm(realm))
}

// DatapointEndpointForRealm returns the endpoint to which datapoints should be
// POSTed for a particular ingest base URL.
func DatapointEndpointForIngestURL(ingestURL string) string {
	return strings.TrimRight(ingestURL, "/") + "/v2/datapoint"
}

// EventEndpointForRealm returns the endpoint to which events should be
// POSTed for a particular realm.
func EventEndpointForRealm(realm string) string {
	return EventEndpointForIngestURL(IngestURLForRealm(realm))
}

// EventEndpointForRealm returns the endpoint to which events should be
// POSTed for a particular ingest base URL.
func EventEndpointForIngestURL(ingestURL string) string {
	return strings.TrimRight(ingestURL, "/") + "/v2/event"
}

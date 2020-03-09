package realm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatapointURLForRealm(t *testing.T) {
	require.Equal(t, DatapointEndpointForRealm("us9"), "https://ingest.us9.signalfx.com/v2/datapoint")
}

func TestEventURLs(t *testing.T) {
	require.Equal(t, EventEndpointForRealm("us9"), "https://ingest.us9.signalfx.com/v2/event")
	require.Equal(t, EventEndpointForIngestURL("https://ingest.us9.signalfx.com/"), "https://ingest.us9.signalfx.com/v2/event")
}

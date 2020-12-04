package messages

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMessage(t *testing.T) {
	t.Run("END_OF_CHANNEL", func(t *testing.T) {
		msg, err := ParseMessage([]byte(`
			{"channel": "ch-1", "event": "END_OF_CHANNEL", "timestampMs": 1607115512410, "type": "control-message"}
		`), true)

		require.NoError(t, err)
		require.IsType(t, &EndOfChannelControlMessage{}, msg)
	})

}

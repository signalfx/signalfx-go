package idtool

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDToString(t *testing.T) {
	assert.Equal(t, "AAAAAAAAAAE", ID(1).String())
	assert.Equal(t, "AAAAAEtEqw4", ID(1262791438).String())
}

func TestIDFromString(t *testing.T) {
	assert.Equal(t, ID(0), IDFromString("AAAAAAAAAAA="))
	assert.Equal(t, ID(0), IDFromString("AAAAAAAAAAA"))
	assert.Equal(t, ID(0), IDFromString("ABCDEFGHIJKKK"))
}

func TestIDUnmarshalJSON(t *testing.T) {
	var id ID
	assert.NoError(t, json.Unmarshal([]byte(`"AAAAAEtEqw4"`), &id))
	assert.Equal(t, ID(1262791438), id)
}
